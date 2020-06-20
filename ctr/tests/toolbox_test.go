package test

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"testing"
	"time"
)

func TestBeegoToolBox(t *testing.T)  {
	tk := toolbox.NewTask("myTask", "0/3 * * * * *", func() error {
		fmt.Println(time.Now())
		return nil
	})

	err := tk.Run()
	if err != nil {
		fmt.Printf("error : %v", err)
	}
	toolbox.AddTask("myTask", tk)
	toolbox.StartTask()

	//for k, v := range toolbox.AdminTaskList {
	//	fmt.Println("key -> ", k, "v -> ")
	//}

	s := "[-1, 80]"
	b := []byte(s)
	var arr []int
	json.Unmarshal(b, &arr)
	fmt.Printf("arr -> %v", arr[0])

	//time.Sleep(1000 * time.Second)
	toolbox.StopTask()
	//toolbox.DeleteTask("myTask")
}