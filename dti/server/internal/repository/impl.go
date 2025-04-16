package repository

type CollegeRepositoryPostgres struct{}

func NewCollegeRepositoryPostgres() CollegeRepository {
	return &CollegeRepositoryPostgres{}
}