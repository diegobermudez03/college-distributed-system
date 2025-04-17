package service

import (
	"errors"

	"github.com/diegobermudez03/college-distributed-system/dti/server/internal/domain"
	"github.com/google/uuid"
)

//poblates the db with the valid faculties and programs, and their relations
func (s *CollegeServiceImpl) PoblateFacultiesAndPrograms() error {
	count, err := s.repository.GetFacultiesCount()
	if err != nil{
		return errors.New("error initializing db")
	}
	if count > 0{
		return nil
	}

	//faculty of ciencias sociales
	socialsFaculty := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Ciencias Sociales",
	}
	socialsFaculty.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Psicologia", FacultyId: socialsFaculty.ID},
		{ID: uuid.New(), Name: "Sociologia", FacultyId: socialsFaculty.ID},
		{ID: uuid.New(), Name: "Trabajo social", FacultyId: socialsFaculty.ID},
		{ID: uuid.New(), Name: "Antropologia", FacultyId: socialsFaculty.ID},
		{ID: uuid.New(), Name: "Comunicacion", FacultyId: socialsFaculty.ID},
	}
	//faculty of ciencias naturles
	scienceFaculty := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Ciencias Naturales",
	}
	scienceFaculty.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Biologia", FacultyId: scienceFaculty.ID},
		{ID: uuid.New(), Name: "Quimica", FacultyId: scienceFaculty.ID},
		{ID: uuid.New(), Name: "Fisica", FacultyId: scienceFaculty.ID},
		{ID: uuid.New(), Name: "Geologia", FacultyId: scienceFaculty.ID},
		{ID: uuid.New(), Name: "Ciencias Ambientales", FacultyId: scienceFaculty.ID},
	}
	//faculty of ingenieria
	engineering := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Ingenieria",
	}
	engineering.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Ingenieria Civil", FacultyId: engineering.ID},
		{ID: uuid.New(), Name: "Ingenieria Electronica", FacultyId: engineering.ID},
		{ID: uuid.New(), Name: "Ingenieria de Sistemas", FacultyId: engineering.ID},
		{ID: uuid.New(), Name: "Ingenieria Mecanica", FacultyId: engineering.ID},
		{ID: uuid.New(), Name: "Ingenieria Industrial", FacultyId: engineering.ID},
	}
	//faculty of medicine
	medicine := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Medicina",
	}
	medicine.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Medicina General", FacultyId: medicine.ID},
		{ID: uuid.New(), Name: "Enfermeria", FacultyId: medicine.ID},
		{ID: uuid.New(), Name: "Odontologia", FacultyId: medicine.ID},
		{ID: uuid.New(), Name: "Farmacia", FacultyId: medicine.ID},
		{ID: uuid.New(), Name: "Terapia Fisica", FacultyId: medicine.ID},
	}
	//faculty of law
	law := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Derecho",
	}
	law.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Derecho Penal", FacultyId: law.ID},
		{ID: uuid.New(), Name: "Derecho Civil", FacultyId: law.ID},
		{ID: uuid.New(), Name: "Derecho Internacional", FacultyId: law.ID},
		{ID: uuid.New(), Name: "Derecho Laboral", FacultyId: law.ID},
		{ID: uuid.New(), Name: "Derecho Constitucional", FacultyId: law.ID},
	}
	//faculty of arts
	arts := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Artes",
	}
	arts.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Bellas Artes", FacultyId: arts.ID},
		{ID: uuid.New(), Name: "Musica", FacultyId: arts.ID},
		{ID: uuid.New(), Name: "Teatro", FacultyId: arts.ID},
		{ID: uuid.New(), Name: "Danza", FacultyId: arts.ID},
		{ID: uuid.New(), Name: "Diseño Grafico", FacultyId: arts.ID},
	}
	//faculty of education
	education := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Educacion",
	}
	education.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Educacion Primaria", FacultyId: education.ID},
		{ID: uuid.New(), Name: "Educacion Secundaria", FacultyId: education.ID},
		{ID: uuid.New(), Name: "Educacion Especial", FacultyId: education.ID},
		{ID: uuid.New(), Name: "Psicopedagogia", FacultyId: education.ID},
		{ID: uuid.New(), Name: "Administracion Educativa", FacultyId: education.ID},
	}
	//faculty of economy
	economy := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Ciencias Economicas",
	}
	economy.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Administracion de Empresas", FacultyId: economy.ID},
		{ID: uuid.New(), Name: "Contabilidad", FacultyId: economy.ID},
		{ID: uuid.New(), Name: "Economia", FacultyId: economy.ID},
		{ID: uuid.New(), Name: "Mercadotecnia", FacultyId: economy.ID},
		{ID: uuid.New(), Name: "Finanzas", FacultyId: economy.ID},
	}
	//faculty of architecture
	architecture := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Arquitectura",
	}
	architecture.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Arquitectura", FacultyId: architecture.ID},
		{ID: uuid.New(), Name: "Urbanismo", FacultyId: architecture.ID},
		{ID: uuid.New(), Name: "Diseño de Interiores", FacultyId: architecture.ID},
		{ID: uuid.New(), Name: "Paisajismo", FacultyId: architecture.ID},
		{ID: uuid.New(), Name: "Restauracion de Patrimonio", FacultyId: architecture.ID},
	}
	//faculty of technology
	tecnology := domain.FacultyModel{
		ID:   uuid.New(),
		Name: "Tecnologia",
	}
	tecnology.Programs = []domain.ProgramModel{
		{ID: uuid.New(), Name: "Desarrollo de Software", FacultyId: tecnology.ID},
		{ID: uuid.New(), Name: "Redes y Telecomunicaciones", FacultyId: tecnology.ID},
		{ID: uuid.New(), Name: "Ciberseguridad", FacultyId: tecnology.ID},
		{ID: uuid.New(), Name: "Inteligencia Artificial", FacultyId: tecnology.ID},
		{ID: uuid.New(), Name: "Big Data", FacultyId: tecnology.ID},
	}

	//save faculties in db
	faculties := []domain.FacultyModel{socialsFaculty, scienceFaculty, engineering, medicine, law, arts, education, economy, architecture, tecnology}
	for _, fac := range faculties{
		if err := s.repository.CreateFaculty(&fac); err != nil{
			return err
		}
	}
	return nil
}
