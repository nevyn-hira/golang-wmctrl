package wmctrl

import(
	"os/exec"
	"bytes"
	"strconv"
	"strings"
)

type Window struct{
	WindowList 	[]*BaseWindow
	xpropPath 	string
	wmctrlPath 	string
}

func New() *Window{
	returnValue := new(Window)
	path, _ := exec.LookPath("xprop")
	returnValue.xpropPath = path
	path, _ = exec.LookPath("wmctrl")
	returnValue.wmctrlPath = path
	return returnValue
}

func (self Window) Get_Active() BaseWindow{
	cmd:=exec.Command(self.xpropPath,"-root","_NET_ACTIVE_WINDOW")
	output, _:=cmd.Output()
	id:=strings.Split(strings.Trim(strings.Split(string(output),"#")[1]," \n"),"x")[1]
	if len(id) < 8{
		for counter:=1;counter<=8-len(id);counter++{
			id="0"+id
		}
	}
	id="0x"+id
	return self.By_ID(id)
}

func (self Window) By_ID(id string) BaseWindow{
	self.List()
	returnValue:=new(BaseWindow)
	for index:=range self.WindowList{
		if self.WindowList[index].ID == id{
			return *self.WindowList[index]
		}
	}
	return *returnValue
}

func (self Window) By_Name(name string) BaseWindow{
	self.List()
	returnValue:=new(BaseWindow)
	for index:=range self.WindowList{
		if self.WindowList[index].WM_Name == name{
			return *self.WindowList[index]
		}
	}
	return *returnValue
}

func (self Window) By_Name_Startswith(prefix string) BaseWindow{
	self.List()
	returnValue:=new(BaseWindow)
	for index:=range self.WindowList{
		if strings.HasPrefix(self.WindowList[index].WM_Name, prefix){
			return *self.WindowList[index]
		}
	}
	return *returnValue
}

func (self Window) By_Name_Endswith(suffix string) BaseWindow{
	self.List()
	returnValue:=new(BaseWindow)
	for index:=range self.WindowList{
		if strings.HasSuffix(self.WindowList[index].WM_Name, suffix){
			return *self.WindowList[index]
		}
	}
	return *returnValue
}

func (self Window) By_Class(class string) BaseWindow{
	self.List()
	returnValue:=new(BaseWindow)
	for index:=range self.WindowList{
		if self.WindowList[index].WM_Class == class{
			return *self.WindowList[index]
		}
	}
	return *returnValue
}

func (self *Window) List(){
	cmd:=exec.Command(self.wmctrlPath,"-lGpx")
	output, _:=cmd.Output()
	test:=bytes.Split(output,[]byte("\n"))
	self.WindowList=self.WindowList[:0]
	for _, data := range test{
		var stringSlice []string
		line := bytes.Split(data,[]byte(" "))
			for _, field:=range(line){
			if len(field)>0{
				stringSlice=append(stringSlice,string(field))
			}
		}
		if len(stringSlice)>0{
			self.WindowList=append(self.WindowList,new(BaseWindow))
			index:=len(self.WindowList)-1
			self.WindowList[index].ID=stringSlice[0]
			strConvert, _:=strconv.Atoi(stringSlice[1])
			self.WindowList[index].Desktop = strConvert
			strConvert, _=strconv.Atoi(stringSlice[2])
			self.WindowList[index].PID = strConvert
			strConvert, _=strconv.Atoi(stringSlice[3])
			self.WindowList[index].X=strConvert
			strConvert, _=strconv.Atoi(stringSlice[4])
			self.WindowList[index].Y=strConvert
			strConvert, _=strconv.Atoi(stringSlice[5])
			self.WindowList[index].H=strConvert
			strConvert, _=strconv.Atoi(stringSlice[6])
			self.WindowList[index].W=strConvert
			self.WindowList[index].WM_Class =stringSlice[7]
			self.WindowList[index].Host = stringSlice[8]
			self.WindowList[index].WM_Name = strings.Join(stringSlice[9:], " ")
			self.WindowList[index].WM_window_role = self.wm_window_role(self.WindowList[index].ID)
			self.WindowList[index].xpropPath = self.xpropPath
			self.WindowList[index].wmctrlPath = self.wmctrlPath
		}
	}
}

func (self Window) wm_window_role(id string) string{
	cmd:=exec.Command(self.xpropPath,"-id",id,"WM_WINDOW_ROLE")
	output, err:=cmd.Output()
	if err!=nil{
		return ""
	} else {
		outputStr:=string(output)
		if strings.Index(outputStr,"=")==-1{
			return ""
		}
		return strings.Trim(strings.Split( outputStr,"=")[1]," \"\n")
	}
}

type BaseWindow struct{
	ID 				string
	Desktop 		int
	PID 			int
	X 				int
	Y 				int
	W 				int
	H 				int
	WM_Class		string
	Host 			string
	WM_Name 		string
	WM_window_role 	string
	xpropPath 		string
	wmctrlPath 		string
}

func (self BaseWindow) IsNull() bool{
	return self.PID==0
}

func (self BaseWindow) Activate(){
	cmd:=exec.Command(self.wmctrlPath,"-id","-a",self.ID)
	cmd.Run()
}

func (self BaseWindow) Resize_and_move(x int,y int,w int,h int){
	cmd:=exec.Command(self.wmctrlPath,"-i","-r",self.ID,"-e","0,"+strconv.Itoa(x)+
		","+strconv.Itoa(y)+","+strconv.Itoa(w)+","+strconv.Itoa(h))
	cmd.Start()
}