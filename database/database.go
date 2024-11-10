package database

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"cloudnotte_practice/graph/model"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var connectionString = "mongodb+srv://Fidel_Wole:2ql24UoUi4uN5302@cluster0.cwzz5uc.mongodb.net/"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	return &DB{client: client}
}

func (db *DB) AddStudent(name string) (*model.Student, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	studentCollec := db.client.Database("school_management_system").Collection("students")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	studentDetails := model.Student{
		Name: name,
	}

	res, err := studentCollec.InsertOne(ctx, studentDetails)
	if err != nil {
		return nil, fmt.Errorf("error adding student: %v", err)
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		studentDetails.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("error converting InsertedID to ObjectID")
	}

	return &studentDetails, nil
}

func (db *DB) AddTerm(name string) (*model.Term, error) {
	if name == "" {
		return nil, fmt.Errorf("term name cannot be empty")
	}

	termCollec := db.client.Database("school_management_system").Collection("terms")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	termDetails := model.Term{
		Name: name,
	}

	res, err := termCollec.InsertOne(ctx, termDetails)
	if err != nil {
		return nil, fmt.Errorf("error adding term: %v", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		termDetails.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("error converting InsertedID to ObjectID")
	}

	return &termDetails, nil
}

func (db *DB) AddSubject(name string) (*model.Subject, error) {
	if name == "" {
		return nil, fmt.Errorf("term name cannot be empty")
	}
	subjectCollec := db.client.Database("school_management_system").Collection("subjects")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	subjectDetails := model.Subject{
		Name: name,
	}

	res, err := subjectCollec.InsertOne(ctx, subjectDetails)
	if err != nil {
		return nil, fmt.Errorf("error adding term: %v", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		subjectDetails.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("error converting InsertedID to ObjectID")
	}
	return &subjectDetails, nil
}

func (db *DB) AddGrade(studentId, termId string, subjects []*model.SubjectGrade) (*model.TermGrade, error) {
	if len(subjects) == 0 {
		return nil, fmt.Errorf("subjects cannot be empty")
	}

	// Convert string IDs to ObjectID
	studentOid, err := primitive.ObjectIDFromHex(studentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student ID: %v", err)
	}

	termOid, err := primitive.ObjectIDFromHex(termId)
	if err != nil {
		return nil, fmt.Errorf("invalid term ID: %v", err)
	}

	student, err := db.GetStudentByID(studentOid)
	if err != nil {
		return nil, fmt.Errorf("error fetching student: %v", err)
	}

	term, err := db.GetTermByID(termOid)
	if err != nil {
		return nil, fmt.Errorf("error fetching term: %v", err)
	}

	var totalMarks float64
	var totalSubjects int
	var subjectPointers []*model.SubjectGrade

	// This Calculates the  total marks and average marks for all subjects
	for _, subject := range subjects {

		weightedMarks := subject.Ca1 +
			subject.Ca2 +
			subject.Obj +
			subject.Theo

		grade := assignGrade(weightedMarks)
		// Store total marks for subject
		subjectGrade := model.SubjectGrade{
			Subject:    &model.Subject{ID: subject.Subject.ID, Name: subject.Subject.Name},
			Ca1:        subject.Ca1,
			Ca2:        subject.Ca2,
			Obj:        subject.Obj,
			Theo:       subject.Theo,
			TotalMarks: weightedMarks,
			Grade:      grade,
		}

		// Add to total marks for the term
		totalMarks += weightedMarks
		totalSubjects++

		// This Append subject grade pointer to the slice
		subjectPointers = append(subjectPointers, &subjectGrade)
	}

	// Calculate average marks for the term
	averageMarks := totalMarks / float64(totalSubjects)

	termGrade := model.TermGrade{
		Student:      &model.Student{ID: studentOid.Hex(), Name: student.Name},
		Term:         &model.Term{ID: termOid.Hex(), Name: term.Name},
		TotalMarks:   totalMarks,
		AverageMarks: averageMarks,
		Subjects:     subjectPointers,
	}

	gradesCollec := db.client.Database("school_management_system").Collection("student_grades")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := gradesCollec.InsertOne(ctx, termGrade)
	if err != nil {
		return nil, fmt.Errorf("error adding grade: %v", err)
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		termGrade.ID = oid.Hex()
	} else {
		return nil, fmt.Errorf("error converting InsertedID to ObjectID")
	}

	return &termGrade, nil
}

func assignGrade(mark float64) string {
	switch {
	case mark >= 70:
		return "A"
	case mark >= 60:
		return "B"
	case mark >= 50:
		return "C"
	case mark >= 40:
		return "D"
	default:
		return "F"
	}
}

func (db *DB) GetSubjectByID(subjectId string) (*model.Subject, error) {
	objectId, err := primitive.ObjectIDFromHex(subjectId)
	if err != nil {
		return nil, fmt.Errorf("invalid subject ID: %v", err)
	}

	subject := &model.Subject{}
	err = db.client.Database("school_management_system").Collection("subjects").FindOne(context.Background(), bson.M{"_id": objectId}).Decode(subject)
	if err != nil {
		return nil, fmt.Errorf("error finding subject: %v", err)
	}

	return subject, nil
}

func (db *DB) GetStudentByID(studentID primitive.ObjectID) (*model.Student, error) {
	student := &model.Student{}
	err := db.client.Database("school_management_system").Collection("students").FindOne(context.Background(), bson.M{"_id": studentID}).Decode(student)
	if err != nil {
		return nil, fmt.Errorf("error finding student: %v", err)
	}
	return student, nil
}

// Fetch Term by ID
func (db *DB) GetTermByID(termID primitive.ObjectID) (*model.Term, error) {
	term := &model.Term{}
	err := db.client.Database("school_management_system").Collection("terms").FindOne(context.Background(), bson.M{"_id": termID}).Decode(term)
	if err != nil {
		return nil, fmt.Errorf("error finding term: %v", err)
	}
	return term, nil
}

func (db *DB) GetStudentGradesByTerm(studentId, termId string) (*model.TermGrade, error) {
	// Convert string IDs to ObjectID
	studentOid, err := primitive.ObjectIDFromHex(studentId)
	if err != nil {
		return nil, fmt.Errorf("invalid student ID: %v", err)
	}

	termOid, err := primitive.ObjectIDFromHex(termId)
	if err != nil {
		return nil, fmt.Errorf("invalid term ID: %v", err)
	}

	filter := bson.M{
		"student.id": studentOid.Hex(),
		"term.id":    termOid.Hex(),
	}

	
	var termGrade model.TermGrade
	err = db.client.Database("school_management_system").Collection("student_grades").
		FindOne(context.Background(), filter).Decode(&termGrade)

	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no grades found for student in the specified term")
		}
		return nil, fmt.Errorf("error fetching student grades: %v", err)
	}

	// For each subject, calculate the position of the student
	for i, subjectGrade := range termGrade.Subjects {
		var allGradesForSubject []model.SubjectGrade
		cursor, err := db.client.Database("school_management_system").
			Collection("student_grades").
			Find(context.Background(), bson.M{
				"term.id":             termOid.Hex(),
				"subjects.subject.id": subjectGrade.Subject.ID,
			})
		if err != nil {
			return nil, fmt.Errorf("error fetching grades for the subject: %v", err)
		}
		defer cursor.Close(context.Background())

		// Collect all grades for the specific subject
		for cursor.Next(context.Background()) {
			var gradeEntry model.TermGrade
			if err := cursor.Decode(&gradeEntry); err != nil {
				return nil, fmt.Errorf("error decoding grade entry: %v", err)
			}

			// Collect subject grades for the matching subject
			for _, sg := range gradeEntry.Subjects {
				if sg.Subject.ID == subjectGrade.Subject.ID {
					allGradesForSubject = append(allGradesForSubject, *sg)
				}
			}
		}

		// This sort the grades by TotalMarks in descending order
		sort.Slice(allGradesForSubject, func(i, j int) bool {
			return allGradesForSubject[i].TotalMarks > allGradesForSubject[j].TotalMarks
		})

		// We go ahead to Find the student's position in the sorted list
		position := -1
		for idx, sg := range allGradesForSubject {
			if sg.Subject.ID == subjectGrade.Subject.ID && sg.TotalMarks == subjectGrade.TotalMarks {
				position = idx + 1
				break
			}
		}

		// We then Add the position to the subject grade
		if position != -1 {
			termGrade.Subjects[i].Position = fmt.Sprintf("Position: %d", position)
		}
	}

	return &termGrade, nil
}
