package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Host struct {
	Id            int       `orm:"column(id);auto" description:"主机id"`
	UserId        *User     `orm:"column(user_id);rel(fk)" description:"用户id"`
	CreateTime    time.Time `orm:"column(create_time);type(timestamp);auto_now_add" description:"创建时间"`
	Ip            string    `orm:"column(ip);size(128)" description:"主机ip"`
	Name          string    `orm:"column(name);size(20);null" description:"主机名"`
	Description   string    `orm:"column(description);size(128);null" description:"主机描述"`
	BaseInfo      string    `orm:"column(base_info);null" description:"主机信息"`
	NeedGetInfo   int8      `orm:"column(need_get_info);null" description:"是否开启监控"`
	GetInfoSpec   string    `orm:"column(get_info_spec);size(20);null" description:"获取主机监控时间间隔"`
	GetInfoTaskId string    `orm:"column(get_info_task_id);size(128);null" description:"主机监控的task id"`
	MemLine       string    `orm:"column(mem_line);size(20);null" description:"内存使用情况范围[0,100]，如果不关心就设置为-1，比如[-1,80]"`
	CpuLine       string    `orm:"column(cpu_line);size(20);null" description:"CPU使用情况范围[0,100]，如果不关心就设置为-1，比如[-1,80]"`
	DiskLine      string    `orm:"column(disk_line);size(20);null" description:"磁盘占用情况范围[0,100]，如果不关心就设置为-1，比如[-1,80]"`
	LoginName     string    `orm:"column(login_name);size(128);null" description:"主机登录用户名"`
	LoginPwd      string    `orm:"column(login_pwd);size(128);null" description:"主机密码"`
	RsaKey        string    `orm:"column(rsa_key);size(1024);null" description:"ssh的public key，暂时不用"`
}

type AddHostReq struct {
	Ip          string `json:"ip"`
	Name        string `json:"name"`
	Description string `json:"description"`
	LoginName   string `json:"login_name"`
	LoginPwd    string `json:"login_pwd"`
}

// 校验添加主机请求参数是否合法
func (req *AddHostReq) Check() (ok bool) {
	if req.Ip == "" || req.Name == "" || req.Description == "" || req.LoginName == "" || req.LoginPwd == "" {
		return false
	}
	return true
}

type HostConnection struct {
	Ip        string `json:"ip"`
	LoginName string `json:"login_name"`
	LoginPwd  string `json:"login_pwd"`
}

type HostProfile struct {
	Id          int           `json:"id"`
	Ip          string        `json:"ip"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	User        UserProfile   `json:"user"` // 创建者
	CreateTime  time.Time     `json:"create_time"`
	BaseInfo    string        `json:"base_info"`
	NeedGetInfo int8          `json:"need_get_info"`
	GetInfoSpec string        `json:"get_info_spec"`
	MemLine     string        `json:"mem_line"`
	CpuLine     string        `json:"cpu_line"`
	DiskLine    string        `json:"disk_line"`
	WatchedUser []UserProfile `json:"watched_user"` // 关注者
	CanWatch    bool          `json:"can_watch"`    // 是否可以关注，如果是关注者则不可以
}

func (h *Host) Host2Profile() (hp HostProfile) {
	var hostWatchs []HostWatch
	o := orm.NewOrm()
	_, _ = o.QueryTable(new(HostWatch)).Filter("host_id", h.Id).All(&hostWatchs)
	hp.CanWatch = true
	hp.WatchedUser = make([]UserProfile, 0)
	for _, hw := range hostWatchs {
		u, _ := GetUserById(hw.UserId.Id)
		up := u.User2UserProfile()
		if up.Id == h.UserId.Id {
			hp.CanWatch = false
		}
		hp.WatchedUser = append(hp.WatchedUser, up)
	}
	hp.Id = h.Id
	hp.Ip = h.Ip
	hp.Name = h.Name
	u, _ := GetUserById(h.UserId.Id)
	hp.User = u.User2UserProfile()
	hp.CreateTime = h.CreateTime
	hp.BaseInfo = h.BaseInfo
	hp.NeedGetInfo = h.NeedGetInfo
	hp.GetInfoSpec = h.GetInfoSpec
	hp.MemLine = h.MemLine
	hp.CpuLine = h.CpuLine
	hp.DiskLine = h.DiskLine
	return hp
}

func (t *Host) TableName() string {
	return "host"
}

func init() {
	orm.RegisterModel(new(Host))
}

// AddHost insert a new Host into database and returns
// last inserted Id on success.
func AddHost(m *Host) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHostById retrieves Host by Id. Returns error if
// Id doesn't exist
func GetHostById(id int) (v *Host, err error) {
	o := orm.NewOrm()
	v = &Host{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHost retrieves all Host matches certain condition. Returns empty list if
// no records exist
func GetAllHost(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Host))
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

	var l []Host
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

// UpdateHost updates Host by Id and returns error if
// the record to be updated doesn't exist
func UpdateHostById(m *Host) (err error) {
	o := orm.NewOrm()
	v := Host{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHost deletes Host by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHost(id int) (err error) {
	o := orm.NewOrm()
	v := Host{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Host{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
