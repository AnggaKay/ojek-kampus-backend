package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/AnggaKay/ojek-kampus-backend/internal/dto"
	"github.com/AnggaKay/ojek-kampus-backend/internal/middleware"
	"github.com/AnggaKay/ojek-kampus-backend/internal/service"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/constants"
	"github.com/AnggaKay/ojek-kampus-backend/pkg/logger"
	"github.com/labstack/echo/v4"
)

type DriverHandler struct {
	driverService service.DriverService
}

func NewDriverHandler(driverService service.DriverService) *DriverHandler {
	return &DriverHandler{
		driverService: driverService,
	}
}

// RegisterDriver handles driver registration with document uploads
// POST /api/auth/register/driver
func (h *DriverHandler) RegisterDriver(c echo.Context) error {
	ctx := c.Request().Context()

	// Parse multipart form with 20MB max
	err := c.Request().ParseMultipartForm(constants.MaxTotalUploadSize)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to parse multipart form")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_FORM", "Failed to parse form data"))
	}

	// Extract form fields
	req := dto.RegisterDriverRequest{
		PhoneNumber:  c.FormValue("phone_number"),
		Password:     c.FormValue("password"),
		FullName:     c.FormValue("full_name"),
		Email:        c.FormValue("email"), // optional
		VehiclePlate: c.FormValue("vehicle_plate"),
		VehicleBrand: c.FormValue("vehicle_brand"),
		VehicleModel: c.FormValue("vehicle_model"),
		VehicleColor: c.FormValue("vehicle_color"),
	}

	// Validate form fields using middleware validator
	if err := middleware.ValidateRequest(c, &req); err != nil {
		return err
	}

	// Extract file uploads
	form := c.Request().MultipartForm
	if form == nil || form.File == nil {
		logger.Log.Warn().Msg("No files uploaded in driver registration")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("MISSING_FILES", "Document files are required"))
	}

	// Helper function to get file from form
	getFile := func(fieldName string) (*multipart.FileHeader, error) {
		files := form.File[fieldName]
		if len(files) == 0 {
			return nil, echo.NewHTTPError(http.StatusBadRequest, fieldName+" document is required")
		}
		return files[0], nil
	}

	// Extract all 4 required documents
	ktpFile, err := getFile("ktp")
	if err != nil {
		logger.Log.Warn().Msg("KTP file missing in driver registration")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("MISSING_KTP", "KTP document is required"))
	}

	simFile, err := getFile("sim")
	if err != nil {
		logger.Log.Warn().Msg("SIM file missing in driver registration")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("MISSING_SIM", "SIM document is required"))
	}

	stnkFile, err := getFile("stnk")
	if err != nil {
		logger.Log.Warn().Msg("STNK file missing in driver registration")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("MISSING_STNK", "STNK document is required"))
	}

	ktmFile, err := getFile("ktm")
	if err != nil {
		logger.Log.Warn().Msg("KTM file missing in driver registration")
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse("MISSING_KTM", "KTM (Student ID Card) document is required"))
	}

	// Prepare documents map for service
	files := map[string]*multipart.FileHeader{
		"ktp":  ktpFile,
		"sim":  simFile,
		"stnk": stnkFile,
		"ktm":  ktmFile,
	}

	// Call service to register driver
	response, err := h.driverService.RegisterDriver(ctx, req, files)
	if err != nil {
		logger.Log.Error().Err(err).Str("phone", req.PhoneNumber).Msg("Driver registration failed")

		// Check error type for appropriate status code
		errMsg := err.Error()
		switch errMsg {
		case constants.ErrPhoneAlreadyRegistered:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse("PHONE_EXISTS", errMsg))
		case constants.ErrEmailAlreadyRegistered:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse("EMAIL_EXISTS", errMsg))
		case constants.ErrVehiclePlateExists:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse("PLATE_EXISTS", errMsg))
		case constants.ErrFileTooLarge:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse("FILE_TOO_LARGE", errMsg))
		case constants.ErrInvalidFileType:
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse("INVALID_FILE_TYPE", errMsg))
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse("REGISTRATION_FAILED", "Failed to register driver"))
		}
	}

	logger.Log.Info().
		Str("phone", req.PhoneNumber).
		Int("driver_id", response.DriverProfile.ID).
		Msg("Driver registered successfully")

	return c.JSON(http.StatusCreated, dto.SuccessResponse("Driver registration successful. Your documents are being reviewed.", response))
}
