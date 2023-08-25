package handler

import (
	"gorest-krsprelist/keluarga"
	"math"
	"strconv"

	// "gorest/faskes"
	"net/http"

	// _ "gorest/docs"

	"github.com/gin-gonic/gin"
)

type keluargaHandler struct {
	keluargaService keluarga.Service
}

func NewKeluargaHandler(keluargaService keluarga.Service) *keluargaHandler {
	return &keluargaHandler{keluargaService}
}

// ListKeluargaSasaran godoc
// @Summary Get Keluarga Sasaran By Kelurahan
// @Description Get Keluarga Sasaran By Kelurahan
// @ID get-list-keluarga-sasaran-by-kelurahan
// @Tags Get Keluarga Sasaran By Kelurahan
// @Accept */*
// @Produce json
// @Success      200  {object}  keluarga.Keluarga
// @Param kdprov query string true "Province Code"
// @Param kdkab query string true "Kabupaten Code"
// @Param kdkec query string true "Kecamatan Code"
// @Param kdkel query string true "Kelurahan Code"
// @Param page query integer false "Page Number"
// @Param pageSize query integer false "Number of Keluarga Sasaran per page" default(10)
// @Router /keluarga [get]
func (handler *keluargaHandler) GetKeluarga(c *gin.Context) {

	KodeDepdagriProvinsi := c.Query("kdprov")
	KodeDepdagriKabupaten := c.Query("kdkab")
	KodeDepdagriKecamatan := c.Query("kdkec")
	KodeDepdagriKelurahan := c.Query("kdkel")

	page := convertToInts(c.DefaultQuery("page", "1"), 1)
	pageSize := convertToInts(c.DefaultQuery("pageSize", "10"), 10)

	if pageSize <= 0 || pageSize > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pageSize value. Must be between 1 and 50",
		})
		return
	}

	keluargas, err := handler.keluargaService.GetKeluarga(KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan, page, pageSize)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors message": err})
		return
	}

	totalItems, err := handler.keluargaService.GetTotalKeluargaCount(KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan)
	if err != nil {
		// Handle the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var keluargasResponse []keluarga.KeluargaResponse

	for _, f := range keluargas {
		keluargaResponse := convertToKeluargaResponse(f)
		keluargasResponse = append(keluargasResponse, keluargaResponse)
	}
	response := gin.H{
		"data": keluargasResponse,
	}

	if totalItems > 0 {
		totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
		response["pagination"] = gin.H{
			"currentPage": page,
			"pageSize":    pageSize,
			"totalPages":  totalPages,
			"totalItems":  totalItems,
		}
	}
	c.JSON(http.StatusOK, response)

}

// ListKrs godoc
// @Summary Get Keluarga Beresiko Stunting By Kelurahan
// @Description Get Keluarga Beresiko Stunting By Kelurahan
// @ID get-list-keluarga-beresiko-stunting-by-kelurahan
// @Tags Get Keluarga Beresiko Stunting By Kelurahan
// @Accept */*
// @Produce json
// @Success      200  {object}  keluarga.Keluarga
// @Param kdprov query string true "Province Code"
// @Param kdkab query string true "Kabupaten Code"
// @Param kdkec query string true "Kecamatan Code"
// @Param kdkel query string true "Kelurahan Code"
// @Param filter2 query integer false "Peringkat Kesejateraan (0-4)"
// @Param page query integer false "Page Number"
// @Param pageSize query integer false "Number of Keluarga Beresiko Stunting per page" default(10)
// @Router /keluargaberesiko [get]
// func (handler *keluargaHandler) GetKeluargaBeresiko(c *gin.Context) {

// 	KodeDepdagriProvinsi := c.Query("kdprov")
// 	KodeDepdagriKabupaten := c.Query("kdkab")
// 	KodeDepdagriKecamatan := c.Query("kdkec")
// 	KodeDepdagriKelurahan := c.Query("kdkel")
// 	// filter1 := c.Query("filter1")
// 	// filter2 := c.Query("filter2")
// 	filter2 := convertToInts(c.Query("filter2"), 0)
// 	page := convertToInts(c.DefaultQuery("page", "1"), 1)
// 	pageSize := convertToInts(c.DefaultQuery("pageSize", "10"), 10)

// 	if pageSize <= 0 || pageSize > 50 {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Invalid pageSize value. Must be between 1 and 50",
// 		})
// 		return
// 	}

// 	keluargas, err := handler.keluargaService.GetKeluargaBeresiko(KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan, filter2, page, pageSize)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"errors message": err})
// 		return
// 	}

// 	totalItems, err := handler.keluargaService.GetTotalKeluargaBeresikoCount(KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan, filter2)
// 	if err != nil {
// 		// Handle the error
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var keluargasResponse []keluarga.KeluargaResponse

// 	for _, f := range keluargas {
// 		keluargaResponse := convertToKeluargaResponse(f)
// 		keluargasResponse = append(keluargasResponse, keluargaResponse)
// 	}
// 	response := gin.H{
// 		"data": keluargasResponse,
// 	}

// 	if totalItems > 0 {
// 		totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
// 		response["pagination"] = gin.H{
// 			"currentPage": page,
// 			"pageSize":    pageSize,
// 			"totalPages":  totalPages,
// 			"totalItems":  totalItems,
// 		}
// 	}
// 	c.JSON(http.StatusOK, response)

// }

func (handler *keluargaHandler) GetKeluargaBeresiko(c *gin.Context) {

	KodeDepdagriProvinsi := c.Query("kdprov")
	KodeDepdagriKabupaten := c.Query("kdkab")
	KodeDepdagriKecamatan := c.Query("kdkec")
	KodeDepdagriKelurahan := c.Query("kdkel")
	// filter1 := c.Query("filter1")
	// filter2 := c.Query("filter2")
	filter2 := convertToInts(c.Query("filter2"), 0)
	page := convertToInts(c.DefaultQuery("page", "1"), 1)
	pageSize := convertToInts(c.DefaultQuery("pageSize", "10"), 10)

	if pageSize <= 0 || pageSize > 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid pageSize value. Must be between 1 and 50",
		})
		return
	}

	keluargas, err := handler.keluargaService.GetKeluargaBeresiko(KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan, filter2, page, pageSize)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors message": err})
		return
	}

	totalItems, err := handler.keluargaService.GetTotalKeluargaBeresikoCount(KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan, filter2)
	if err != nil {
		// Handle the error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var keluargasResponse []keluarga.KeluargaResponse

	for _, f := range keluargas {
		keluargaResponse := convertToKeluargaResponse(f)
		keluargasResponse = append(keluargasResponse, keluargaResponse)
	}
	response := gin.H{
		"data": keluargasResponse,
	}

	if totalItems > 0 {
		totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))
		response["pagination"] = gin.H{
			"currentPage": page,
			"pageSize":    pageSize,
			"totalPages":  totalPages,
			"totalItems":  totalItems,
		}
	}
	c.JSON(http.StatusOK, response)

}

func (handler *keluargaHandler) SearchByNik(c *gin.Context) {

	nik := c.Query("nik")
	filter2 := convertToInts(c.Query("filter2"), 0)
	filter3 := c.Query("filter3")
	filter4 := c.Query("filter4")

	result, err := handler.keluargaService.SearchByNik(nik, filter2, filter3, filter4)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": err.Error(),
		})
		return
	}

	if result == nil {
		// Handle the case where no records were found
		errorMessage := "No records found"
		c.JSON(http.StatusNotFound, gin.H{
			"error_message": errorMessage,
		})
		return
	}

	var keluargasResponse []keluarga.KeluargaResponse

	for _, f := range result {
		keluargaResponse := convertToKeluargaResponse(f)
		keluargasResponse = append(keluargasResponse, keluargaResponse)
	}

	response := gin.H{
		"data": keluargasResponse,
	}

	c.JSON(http.StatusOK, response)
}

func convertToInts(str string, defaultValue int) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return val
}

func convertToKeluargaResponse(f keluarga.Keluarga) keluarga.KeluargaResponse {
	return keluarga.KeluargaResponse{

		// IdKabupaten:                         int(f.IdKabupaten),
		// KodeKabupaten:                       f.KodeKabupaten,
		// NamaKabupaten:                       f.NamaKabupaten,
		// IdProvinsi:                          int(f.IdProvinsi),
		// KodeProvinsi:                        f.KodeProvinsi,
		// NamaProvinsi:                        f.NamaProvinsi,
		// JumlahKeluarga:                      int(f.JumlahKeluarga),
		// JumlahKeluargaSasaran:               int(f.JumlahKeluargaSasaran),
		// Prioritas1:                          int(f.Prioritas1),
		// Prioritas2:                          int(f.Prioritas2),
		// Prioritas3:                          int(f.Prioritas3),
		// Prioritas4:                          int(f.Prioritas4),
		// PeringkatKesejahteraanDiatas4:       int(f.PeringkatKesejahteraanDiatas4),
		// JumlahKeluargaBeresikoStunting:      int(f.JumlahKeluargaBeresikoStunting),
		// JumlahPus:                           int(f.JumlahPus),
		// JumlahPusHamil:                      int(f.JumlahPusHamil),
		// JumlahBalita:                        int(f.JumlahBalita),
		// JumlahBaduta:                        int(f.JumlahBaduta),
		// JumlahBukanPesertaKbModern:          int(f.JumlahBukanPesertaKbModern),
		// JumlahJambanTidakLayak:              int(f.JumlahJambanTidakLayak),
		// JumlahAirTidakLayak:                 int(f.JumlahAirTidakLayak),
		// JumlahKeluargaTidakBeresikoStunting: int(f.JumlahKeluargaTidakBeresikoStunting),

		IdFrm:                  f.IdFrm,
		IdProvinsi:             int(f.IdProvinsi),
		KodeProvinsi:           int(f.KodeProvinsi),
		NamaProvinsi:           f.NamaProvinsi,
		IdKabupaten:            int(f.IdKabupaten),
		KodeKabupaten:          int(f.KodeKabupaten),
		NamaKabupaten:          f.NamaKabupaten,
		IdKecamatan:            int(f.IdKecamatan),
		KodeKecamatan:          int(f.KodeKecamatan),
		NamaKecamatan:          f.NamaKecamatan,
		IdKelurahan:            int(f.IdKelurahan),
		KodeKelurahan:          int(f.KodeKelurahan),
		NamaKelurahan:          f.NamaKelurahan,
		KodeRw:                 f.KodeRw,
		NamaRw:                 f.NamaRw,
		KodeRt:                 f.KodeRt,
		NamaRt:                 f.NamaRt,
		Periode:                int(f.Periode),
		KodeKeluarga:           f.KodeKeluarga,
		Nik:                    f.Nik,
		NamaKepalaKeluarga:     f.NamaKepalaKeluarga,
		NamaIstri:              f.NamaIstri,
		Baduta:                 f.Baduta,
		Balita:                 f.Balita,
		Pus:                    f.Pus,
		PusHamil:               f.PusHamil,
		SumberAirLayakTidak:    f.SumberAirLayakTidak,
		JambanLayakTidak:       f.JambanLayakTidak,
		TerlaluMuda:            f.TerlaluMuda,
		TerlaluTua:             f.TerlaluTua,
		TerlaluDekat:           f.TerlaluDekat,
		TerlaluBanyak:          f.TerlaluBanyak,
		BukanPesertaKbModern:   f.BukanPesertaKbModern,
		KesejahteraanPrioritas: int(f.KesejahteraanPrioritas),
		RisikoStunting:         f.RisikoStunting,
		Latitude:               f.Latitude,
		Longitude:              f.Longitude,
	}
}
