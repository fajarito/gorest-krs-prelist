package keluarga

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	// FindByID(nik string) (Stunting, error)
	// FindByPusHamil(nik string) (Stunting, error)
	// FindByStunting(nik string) (Stunting, error)
	// FindByBaduta(nik string) (Stunting, error)
	//KrsByKec(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string) ([]Krk, error)
	GetKeluarga(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string, page int, pageSize int) ([]Keluarga, error)
	GetTotalKeluargaCount(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string) (int64, error)
	GetKeluargaBeresiko(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string, filter2 int, page int, pageSize int) ([]Keluarga, error)
	GetTotalKeluargaBeresikoCount(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string, filter2 int) (int64, error)
	SearchByNik(nik string, filter2 int, filter3 string, filter4 string) ([]Keluarga, error)

	// FindByKec(ProvinsiID int, KabupatenID int, KecamatanID int) ([]Krs, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetKeluarga(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string, page int, pageSize int) ([]Keluarga, error) {
	var keluargas []Keluarga

	// err := r.db.Find(&faskeses, ProvinsiID).Error
	offset := (page - 1) * pageSize
	err := r.db.Where("kode_depdagri_provinsi = ? AND kode_depdagri_kabupaten = ? AND kode_depdagri_kecamatan = ? AND kode_depdagri_kelurahan = ?", KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan).Offset(offset).Limit(pageSize).Find(&keluargas).Error

	if err != nil {
		return nil, err
	}

	return keluargas, nil
}

func (r *repository) GetTotalKeluargaCount(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string) (int64, error) {
	var count int64
	if err := r.db.Model(&Keluarga{}).Where("kode_depdagri_provinsi = ? AND kode_depdagri_kabupaten = ? AND kode_depdagri_kecamatan = ? AND kode_depdagri_kelurahan = ?", KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) GetKeluargaBeresiko(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string, filter2 int, page int, pageSize int) ([]Keluarga, error) {
	var keluargas []Keluarga

	offset := (page - 1) * pageSize
	query := r.db.Where("kode_depdagri_provinsi = ? AND kode_depdagri_kabupaten = ? AND kode_depdagri_kecamatan = ? AND kode_depdagri_kelurahan = ? AND risiko_stunting= 'V'", KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan)

	if filter2 != 0 {
		allowedFilter2Values := []int{1, 2, 3, 4, 0}
		isFilter2Valid := false

		for _, value := range allowedFilter2Values {
			if filter2 == value {
				isFilter2Valid = true
				break
			}
		}

		if !isFilter2Valid {
			return nil, errors.New("Invalid filter2 value. Must be one of: 1, 2, 3, 4, 0")
		}

		query = query.Where("kesejahteraan_prioritas = ?", filter2)
	}

	err := query.Offset(offset).Limit(pageSize).Find(&keluargas).Error
	if err != nil {
		return nil, err
	}

	return keluargas, nil
}

// func (r *repository) GetTotalKeluargaBeresikoCount(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string) (int64, error) {
// 	var count int64
// 	if err := r.db.Model(&Keluarga{}).Where("kode_depdagri_provinsi = ? AND kode_depdagri_kabupaten = ? AND kode_depdagri_kecamatan = ? AND kode_depdagri_kelurahan = ? AND risiko_stunting= 'V'", KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan).Count(&count).Error; err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

func (r *repository) GetTotalKeluargaBeresikoCount(KodeDepdagriProvinsi string, KodeDepdagriKabupaten string, KodeDepdagriKecamatan string, KodeDepdagriKelurahan string, filter2 int) (int64, error) {
	var count int64
	query := r.db.Model(&Keluarga{}).Where("kode_depdagri_provinsi = ? AND kode_depdagri_kabupaten = ? AND kode_depdagri_kecamatan = ? AND kode_depdagri_kelurahan = ? AND risiko_stunting = 'V'", KodeDepdagriProvinsi, KodeDepdagriKabupaten, KodeDepdagriKecamatan, KodeDepdagriKelurahan)

	if filter2 != 0 {
		query = query.Where("kesejahteraan_prioritas = ?", filter2)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *repository) SearchByNik(nik string, filter2 int, filter3 string, filter4 string) ([]Keluarga, error) {
	var keluargas []Keluarga

	query := r.db.Where("nik = ?", nik)

	/*
		Filter 2 = Desil (0,1,2,3,4)
		Filter 3 = Pus Hamil (V,X)
		Filter 4 = Risiko Stunting (V,X)
	*/
	if filter2 != 0 {
		allowedFilter2Values := []int{1, 2, 3, 4, 0}
		isFilter2Valid := false

		for _, value := range allowedFilter2Values {
			if filter2 == value {
				isFilter2Valid = true
				break
			}
		}

		if !isFilter2Valid {
			return nil, errors.New("Invalid filter2 value. Must be one of: 1, 2, 3, 4, 0")
		}

		query = query.Where("kesejahteraan_prioritas = ?", filter2)
	}

	if filter3 == "V" || filter3 == "X" {
		query = query.Where("pus_hamil = ?", filter3)
	} else if filter3 != "" {
		return nil, errors.New("Invalid filter3 value. Must be 'V' or 'X'")
	}

	if filter4 == "V" || filter4 == "X" {
		query = query.Where("risiko_stunting = ?", filter4)
	} else if filter4 != "" {
		return nil, errors.New("Invalid filter3 value. Must be 'V' or 'X'")
	}

	err := query.Find(&keluargas).Error
	if err != nil {
		// Handle the error and send an error response
		errorMessage := "An error occurred while fetching data"
		// You can customize the error message based on the specific error if needed
		return nil, errors.New(errorMessage)
	}

	if len(keluargas) == 0 {
		// Send a custom response indicating no records found
		return nil, errors.New("No records found")
	}

	return keluargas, nil
}
