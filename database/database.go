package database

import (
	"context"
	"fmt"
	"log"
	//"sort"
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
		fmt.Printf("Subject ID: %s, Subject Name: %s\n", subject.Subject.ID, subject.Subject.Name)
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
	// Convert string ID to MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(subjectId)
	if err != nil {
		return nil, fmt.Errorf("invalid subject ID: %v", err)
	}

	// Temporary struct to map MongoDB `_id` to Go's ObjectID type
	var mongoSubject struct {
		ID   primitive.ObjectID `bson:"_id"`
		Name string             `bson:"name"`
	}

	// Retrieve the subject from MongoDB
	err = db.client.Database("school_management_system").
		Collection("subjects").
		FindOne(context.Background(), bson.M{"_id": objectId}).
		Decode(&mongoSubject)

	if err != nil {
		return nil, fmt.Errorf("error finding subject: %v", err)
	}

	// Convert MongoDB ObjectID to string for the gqlgen Subject struct
	subject := &model.Subject{
		ID:   mongoSubject.ID.Hex(),
		Name: mongoSubject.Name,
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
    // This Convert string IDs to ObjectID
    studentOid, err := primitive.ObjectIDFromHex(studentId)
    if err != nil {
        return nil, fmt.Errorf("invalid student ID: %v", err)
    }

    termOid, err := primitive.ObjectIDFromHex(termId)
    if err != nil {
        return nil, fmt.Errorf("invalid term ID: %v", err)
    }

    // Filter to find the specific student's grades for the term
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

    // For each subject, calculate the student's position using MongoDB aggregation
    for i, subjectGrade := range termGrade.Subjects {
        fmt.Printf("Processing subject ID: %s\n", subjectGrade.Subject.ID)

        pipeline := mongo.Pipeline{
            // Match documents for the specific term
            bson.D{
                {Key: "$match", Value: bson.D{
                    {Key: "term.id", Value: termOid.Hex()},
                }},
            },
            // Unwind subjects to get individual subject records
            bson.D{
                {Key: "$unwind", Value: "$subjects"},
            },
            // Match the specific subject we're interested in
            bson.D{
                {Key: "$match", Value: bson.D{
                    {Key: "subjects.subject.id", Value: subjectGrade.Subject.ID},
                }},
            },
            // Project the relevant fields
            bson.D{
                {Key: "$project", Value: bson.D{
                    {Key: "studentId", Value: "$student.id"},
                    {Key: "totalMarks", Value: "$subjects.totalmarks"}, // Fixed field name to match document structure
                }},
            },
            // Sort by total marks in descending order
            bson.D{
                {Key: "$sort", Value: bson.D{
                    {Key: "totalMarks", Value: -1},
                }},
            },
            // Add dense rank
            bson.D{
                {Key: "$setWindowFields", Value: bson.D{
                    {Key: "partitionBy", Value: bson.D{}},
                    {Key: "sortBy", Value: bson.D{
                        {Key: "totalMarks", Value: -1},
                    }},
                    {Key: "output", Value: bson.D{
                        {Key: "position", Value: bson.D{
                            {Key: "$denseRank", Value: bson.D{}},
                        }},
                    }},
                }},
            },
            // This Match only the current student to get their position
            bson.D{
                {Key: "$match", Value: bson.D{
                    {Key: "studentId", Value: studentOid.Hex()},
                }},
            },
        }

        // Debug: Print the pipeline
        fmt.Printf("Aggregation Pipeline for subject %s: %+v\n", subjectGrade.Subject.ID, pipeline)

        cursor, err := db.client.Database("school_management_system").
            Collection("student_grades").
            Aggregate(context.Background(), pipeline)
        if err != nil {
            return nil, fmt.Errorf("error calculating position for subject: %v", err)
        }
        defer cursor.Close(context.Background())

        var allResults []struct {
            StudentID  string `bson:"studentId"`
            Position   int    `bson:"position"`
            TotalMarks int    `bson:"totalMarks"`
        }

        if err := cursor.All(context.Background(), &allResults); err != nil {
            return nil, fmt.Errorf("error reading all results: %v", err)
        }

        fmt.Printf("Results for subject %s: %+v\n", subjectGrade.Subject.ID, allResults)

        if len(allResults) > 0 {
            result := allResults[0]
            fmt.Printf("Found position for student %s: Position=%d, TotalMarks=%d\n",
                result.StudentID, result.Position, result.TotalMarks)
            termGrade.Subjects[i].Position = fmt.Sprintf("Position: %d", result.Position)
        } else {
            fmt.Printf("No position found for student %s in subject %s\n",
                studentOid.Hex(), subjectGrade.Subject.ID)
            termGrade.Subjects[i].Position = "Position not available"
        }
    }

    return &termGrade, nil
}