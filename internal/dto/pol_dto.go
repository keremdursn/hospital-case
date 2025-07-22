package dto

type PolyclinicLookup struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AddHospitalPolyclinicRequest struct {
	PolyclinicID uint `json:"polyclinic_id"`
}

type HospitalPolyclinicResponse struct {
	ID             uint   `json:"id"`
	PolyclinicID   uint   `json:"polyclinic_id"`
	PolyclinicName string `json:"polyclinic_name"`
}
