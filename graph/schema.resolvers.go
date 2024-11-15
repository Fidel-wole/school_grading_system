package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.56

import (
	"cloudnotte_practice/database"
	"cloudnotte_practice/graph/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)
var db = database.Connect()
// AddStudent is the resolver for the addStudent field.
func (r *mutationResolver) AddStudent(ctx context.Context, name string) (*model.Student, error) {
	return db.AddStudent(name)
}

// AddSubject is the resolver for the addSubject field.
func (r *mutationResolver) AddSubject(ctx context.Context, name string) (*model.Subject, error) {
	return db.AddSubject(name)
}

// AddTerm is the resolver for the addTerm field.
func (r *mutationResolver) AddTerm(ctx context.Context, name string) (*model.Term, error) {
	return db.AddTerm(name)
}

// AddGrade is the resolver for the addGrade field.
func (r *mutationResolver) AddGrade(ctx context.Context, studentID string, termID string, subjects []*model.SubjectGradeInput) (*model.TermGrade, error) {
	var subjectGrades []*model.SubjectGrade
	for _, subjectInput := range subjects {
		subject, err := db.GetSubjectByID(subjectInput.SubjectID)
		if err != nil {
			return nil, fmt.Errorf("subject not found: %v", err)
		}
		fmt.Printf("Subject ID: %v\n", subject.ID)
		// This Creates a new SubjectGrade instance based on SubjectGradeInput
		subjectGrade := &model.SubjectGrade{
			Subject: subject,
			Ca1:     subjectInput.Ca1,
			Ca2:     subjectInput.Ca2,
			Obj:     subjectInput.Obj,
			Theo:    subjectInput.Theo,
		}
		subjectGrades = append(subjectGrades, subjectGrade)
	}

	return db.AddGrade(studentID, termID, subjectGrades)
}

// GetStudent is the resolver for the getStudent field.
func (r *queryResolver) GetStudent(ctx context.Context, id string) (*model.Student, error) {
	// Convert the string ID to ObjectID
	studentID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid student ID format: %v", err)
	}

	return db.GetStudentByID(studentID)
}

// GetAllStudents is the resolver for the getAllStudents field.
func (r *queryResolver) GetAllStudents(ctx context.Context) ([]*model.Student, error) {
	panic(fmt.Errorf("not implemented: GetAllStudents - getAllStudents"))
}

// GetStudentGradesByTerm is the resolver for the getStudentGradesByTerm field.
func (r *queryResolver) GetStudentGradesByTerm(ctx context.Context, studentID string, termID string) (*model.TermGrade, error) {
	return db.GetStudentGradesByTerm(studentID, termID)
}

// GetStudentCumulativeGrades is the resolver for the getStudentCumulativeGrades field.
func (r *queryResolver) GetStudentCumulativeGrades(ctx context.Context, studentID string) ([]*model.CumulativeGrade, error) {
	panic(fmt.Errorf("not implemented: GetStudentCumulativeGrades - getStudentCumulativeGrades"))
}

// GetSubjectByID is the resolver for the getSubjectById field.
func (r *queryResolver) GetSubjectByID(ctx context.Context, subjectID string) (*model.Subject, error) {
	panic(fmt.Errorf("not implemented: GetSubjectByID - getSubjectById"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	var db = database.Connect()
*/
