package task_pages

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const (
	WeightSetsAndReps int = iota
	TimePaceAndDistance
)

type TaskTable struct {
	Headers []string
	Rows    []TaskRow
}

type TaskRow struct {
	Value []string
}

type TaskModel struct {
	Title     string
	Date      time.Time
	Value     string
	Type      int
	Id        int
	TaskTable TaskTable
}

type AddTaskModel struct {
	Title string
	Value string
	Type  int
}

type WeightSetsAndRepsModel struct {
	Weight int
	Sets   int
	Reps   int
}

type TimePaceAndDistanceModel struct {
	Pace string
	Time string
}

func (model WeightSetsAndRepsModel) ToTable() TaskTable {
	table := TaskTable{
		Headers: []string{"Weight", "Reps", "Sets"},
		Rows: []TaskRow{
			{[]string{strconv.Itoa(model.Weight), strconv.Itoa(model.Reps), strconv.Itoa(model.Sets)}},
		},
	}
	return table
}

func (model TimePaceAndDistanceModel) ToTable() TaskTable {
	table := TaskTable{
		Headers: []string{"Pace", "Time"},
		Rows: []TaskRow{
			{[]string{model.Pace, model.Time}},
		},
	}
	return table
}

func WeightSetsAndRepsFromModel(value string) WeightSetsAndRepsModel {
	model := WeightSetsAndRepsModel{}
	err := json.Unmarshal([]byte(value), &model)
	if err != nil {
		return WeightSetsAndRepsModel{}
	}
	return model
}

func TimePaceAndDistanceFromModel(value string) TimePaceAndDistanceModel {
	model := TimePaceAndDistanceModel{}
	err := json.Unmarshal([]byte(value), &model)
	if err != nil {
		return TimePaceAndDistanceModel{}
	}
	return model
}

func WeightSetsAndRepsFromRequest(r *http.Request) (WeightSetsAndRepsModel, error) {
	weight, err := strconv.Atoi(r.FormValue("weight"))
	if err != nil {
		return WeightSetsAndRepsModel{}, err
	}

	sets, err := strconv.Atoi(r.FormValue("sets"))
	if err != nil {
		return WeightSetsAndRepsModel{Weight: weight}, err
	}

	reps, err := strconv.Atoi(r.FormValue("reps"))
	if err != nil {
		return WeightSetsAndRepsModel{Weight: weight, Sets: sets}, err
	}

	return WeightSetsAndRepsModel{
		Weight: weight,
		Sets:   sets,
		Reps:   reps,
	}, nil
}

func TimePaceAndDistanceFromRequest(r *http.Request) TimePaceAndDistanceModel {
	return TimePaceAndDistanceModel{
		Pace: r.FormValue("pace"),
		Time: r.FormValue("time"),
	}
}

func AddTaskModelFromRequest(r *http.Request) (AddTaskModel, error) {
	typeId, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		return AddTaskModel{}, err
	}

	var value string
	switch typeId {
	case WeightSetsAndReps:
		model, err := WeightSetsAndRepsFromRequest(r)
		if err != nil {
			return AddTaskModel{}, err
		}

		valueBytes, err := json.Marshal(model)
		if err != nil {
			return AddTaskModel{}, err
		}
		value = string(valueBytes)
	case TimePaceAndDistance:
		model := TimePaceAndDistanceFromRequest(r)

		valueBytes, err := json.Marshal(model)
		if err != nil {
			return AddTaskModel{}, err
		}
		value = string(valueBytes)
	}

	return AddTaskModel{
		Title: r.FormValue("title"),
		Value: value,
		Type:  typeId,
	}, nil
}
