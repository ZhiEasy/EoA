package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Task struct {
	Id          int       `orm:"column(id);auto" description:"任务id"`
	UserId      *User     `orm:"column(user_id);rel(fk)" description:"用户id"`
	HostId      *Host     `orm:"column(host_id);rel(fk)" description:"主机id"`
	Name        string    `orm:"column(name);size(128)" description:"任务名称"`
	Description string    `orm:"column(description);size(128)" description:"任务描述"`
	Spec        string    `orm:"column(spec);size(20)" description:"任务执行时间"`
	Type        int       `orm:"column(type)" description:"任务类型"`
	CreateTime  time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
}

type TaskProfile struct {
	Id          int         `json:"id"`
	User        UserProfile `json:"user"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Spec        string      `json:"spec"`
	Type        int         `json:"type"`
	CreateTime  time.Time   `json:"create_time"`
}

func (t *Task) Task2Profile() (tp *TaskProfile) {
	userObj, err := GetUserById(t.UserId.Id)
	if err != nil {
		return nil
	}
	return &TaskProfile{
		Id:          t.Id,
		User:        userObj.User2UserProfile(),
		Name:        t.Name,
		Description: t.Description,
		Spec:        t.Spec,
		Type:        t.Type,
		CreateTime:  t.CreateTime,
	}
}

type AddHostInfoTask struct {
	HostId      int       `json:"host_id"`
	Description string    `json:"description"`
	Spec        string    `json:"spec"`
	MemLine     []float64 `json:"mem_line"`
	CpuLine     []float64 `json:"cpu_line"`
	DiskLine    []float64 `json:"disk_line"`
}

func (t *Task) TableName() string {
	return "task"
}

func init() {
	orm.RegisterModel(new(Task))
}

// AddTask insert a new Task into database and returns
// last inserted Id on success.
func AddTask(m *Task) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTaskById retrieves Task by Id. Returns error if
// Id doesn't exist
func GetTaskById(id int) (v *Task, err error) {
	o := orm.NewOrm()
	v = &Task{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllTask retrieves all Task matches certain condition. Returns empty list if
// no records exist
func GetAllTask(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Task))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Task
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateTask updates Task by Id and returns error if
// the record to be updated doesn't exist
func UpdateTaskById(m *Task) (err error) {
	o := orm.NewOrm()
	v := Task{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTask deletes Task by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTask(id int) (err error) {
	o := orm.NewOrm()
	v := Task{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Task{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
