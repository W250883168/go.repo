'use strict';
/**
 * Created by Administrator on 2016/9/8.
 */

app.filter('propsFilter', function() {
	return function(items, props) {
		var out = [];
		if(angular.isArray(items)) {
			items.forEach(function(item) {
				var itemMatches = false;
				var keys = Object.keys(props);
				for(var i = 0; i < keys.length; i++) {
					var prop = keys[i];
					try {
						if(typeof(props[prop]) == "number") {
							var text = props[prop];
							if(item[prop].toString() == text.toString()) {
								itemMatches = true;
								break;
							}
						} else {
							var text = props[prop].toString().toLowerCase();
							if(item[prop].toString().toLowerCase().indexOf(text) !== -1) {
								itemMatches = true;
								break;
							}
						}
					} catch(e) {}
				}
				if(itemMatches) {
					out.push(item);
				}
			});
		} else {
			out = items;
		}
		return out;
	};
});

/*  ------------  弹窗  --------------  */

//    弹窗1
app.controller('modalGetClassRoomCtrl', ['$scope', '$modalInstance', 'items', 'httpService', function($scope, $modalInstance, items, httpService) {
	//$scope.items = items;

	/*---------------- 取校区 楼栋 楼层 教室 -----------------*/
	//  校区
	$scope.school = {};
	$scope.school.schoolList = "";
	$scope.school.schoolItems = [];
	//  楼栋
	$scope.building = {};
	$scope.building.buildingList = "";
	$scope.building.buildingItems = [];
	//  楼层
	$scope.floors = {};
	$scope.floors.floorsList = "";
	$scope.floors.floorsItems = [];
	//  教室
	$scope.classroom = {};
	$scope.classroom.classroomList = "";
	$scope.classroom.classroomItems = [];
	//   取校区列表
	$scope.getcampus = function() {
		var url = config.HttpUrl + "/basicset/getcampus";
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.school.schoolItems = data.Result;
				//console.log($scope.school.schoolItems)
			} else {
				alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	$scope.getcampus();

	//  取楼栋列表
	$scope.floorItems = [];
	$scope.getcampusFloor = function(campusid) {
			if(!campusid) return;
			var url = config.HttpUrl + "/basicset/getbuilding";
			var data = {
				"campusid": Number(campusid)
			};
			var promise = httpService.ajaxGet(url, data);
			promise.then(function(data) {
				//console.log("取楼栋")
				if(data.Rcode == "1000") {
					$scope.building.buildingItems = data.Result;
					//console.log($scope.building.buildingItems)
				} else {
					$scope.tipService.setMessage(data.Reason);
					//alert(data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//   取楼层和教室列表
	$scope.getfloorsandrooms = function(buildingid) {
			var url = config.HttpUrl + "/basicset/getfloorsandrooms";
			var data = {
				"buildingid": buildingid
			};
			var promise = httpService.ajaxGet(url, data);
			promise.then(function(data) {
				if(data.Rcode == "1000") {
					$scope.floors.floorsItems = data.Result;
					//console.log($scope.floors.floorsItems)
				} else {
					$scope.tipService.setMessage(data.Reason);
					//alert(data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//   取教室列表
	$scope.getClassRooms = function(floorid) {
		if($scope.floors.floorsItems.length > 0) {
			for(var i = 0; i < $scope.floors.floorsItems.length; i++) {
				if($scope.floors.floorsItems[i].Floorsid == floorid) {
					$scope.classroom.classroomItems = $scope.floors.floorsItems[i].Rooms;
				}
			}
		}
	}

	/*---------------- 取校区 楼栋 楼层 教室 End -----------------*/

	$scope.changeSelect = function(types) {
		switch(types) {
			case "school":
				$scope.building.buildingList = "";
				$scope.building.buildingItems = [];
				$scope.floors.floorsList = "";
				$scope.floors.floorsItems = [];
				$scope.classroom.classroomList = "";
				$scope.classroom.classroomItems = [];
				$scope.getcampusFloor($scope.school.schoolList);
				break;
			case "building":
				$scope.floors.floorsList = "";
				$scope.floors.floorsItems = [];
				$scope.classroom.classroomList = "";
				$scope.classroom.classroomItems = [];
				$scope.getfloorsandrooms($scope.building.buildingList);
				break;
			case "floor":
				$scope.classroom.classroomList = "";
				$scope.getClassRooms($scope.floors.floorsList);

				break;
			case "room":

				break;
			default:
				;
		}
	}

	$scope.ok = function() {
		var list = "";
		var add = "";
		if($scope.school.schoolList != "") {
			for(var i = 0; i < $scope.school.schoolItems.length; i++) {
				if($scope.school.schoolList == $scope.school.schoolItems[i].Campusid) {
					add += $scope.school.schoolItems[i].Campusname;
				}
			}
			list = {
				"add": add,
				"addId": $scope.school.schoolList,
				"addCode": "campus"
			};
		}
		if($scope.building.buildingList != "") {
			for(var i = 0; i < $scope.building.buildingItems.length; i++) {
				if($scope.building.buildingList == $scope.building.buildingItems[i].Buildingid) {
					add += "-" + $scope.building.buildingItems[i].Buildingname;
				}
			}
			list = {
				"add": add,
				"addId": $scope.building.buildingList,
				"addCode": "building"
			};
		}
		if($scope.floors.floorsList != "") {
			for(var i = 0; i < $scope.floors.floorsItems.length; i++) {
				if($scope.floors.floorsList == $scope.floors.floorsItems[i].Floorsid) {
					add += "-" + $scope.floors.floorsItems[i].Floorname;
				}
			}
			list = {
				"add": add,
				"addId": $scope.floors.floorsList,
				"addCode": "floor"
			};
		}
		if($scope.classroom.classroomList != "") {
			for(var i = 0; i < $scope.classroom.classroomItems.length; i++) {
				if($scope.classroom.classroomList == $scope.classroom.classroomItems[i].Classroomid) {
					add += "-" + $scope.classroom.classroomItems[i].Classroomsname;
				}
			}
			list = {
				"add": add,
				"addId": $scope.classroom.classroomList,
				"addCode": "classroom"
			};
		}
		$modalInstance.close(list);
	};

	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};
}]);

//    选择设备型号弹窗
app.controller("modalGetDeviceCtrl", ['$scope', 'httpService', '$modalInstance','items', function($scope, httpService, $modalInstance,items) {
	//console.log("选择设备型号弹窗");
	
	//   items:{'Type':'2','show':false}   用于是否可选
	$scope.items = items;
	
	//   设置哪些可先    
	//   return obj
	//   showFilte==true  可选   false 不可选
	//   items:{'Type':'2','show':false}   用于是否可选
	var showfilte = function(items,deviceitems){
		if(!items || !deviceitems){
			for(var a in deviceitems){
				deviceitems[a].showFilte = true;
			}
			return deviceitems;
		}
		for(var a in deviceitems){
			if(deviceitems[a].Type == items.Type){
				deviceitems[a].showFilte = items.show;
			}else{
				deviceitems[a].showFilte = !items.show;
			}
		}
		return deviceitems;
	}

	$scope.xhgl_list = [];
	var init_data = function() {
		var url = config.HttpUrl + "/device/getDeviceModelTree";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.xhgl_list = data.Result.Data;
				$scope.xhgl_list = showfilte(items,$scope.xhgl_list);
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log(data)
		}, function(reason) {}, function(update) {});
	};
	// run
	init_data();

	//  查找叶子
	$scope.isLeaf = function(id) {
		var bol = true;
		for(var i = 0; i < $scope.xhgl_list.length; i++) {
			if($scope.xhgl_list[i].PId == id) {
				bol = false;
				break;
			}
		}
		return bol;
	}

	//   toggle
	$scope.isOpen = function(item) {
		return item.bul = !item.bul;
	}

	$scope.ok = function(item) {
		//console.log(item)
		$modalInstance.close(item);
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};

}]);

/*   弹窗  -  选择老师       */
app.controller("modalGetTeacherCtrl", ['$scope', 'httpService', '$modalInstance', function($scope, httpService, $modalInstance) {
	//console.log("弹窗-选择老师");

	$scope.getall_data = [];
	//  学院
	$scope.college = {};
	$scope.college.collegeList = "";
	$scope.college.collegeItems = [];
	//   科系
	$scope.major = {};
	$scope.major.majorList = "";
	$scope.major.majorItems = [];
	//   老师
	$scope.teacher = {};
	$scope.teacher.teacherList = "";
	$scope.teacher.teacherItems = [];

	//   取校学院列表
	$scope.getall = function() {
		var url = config.HttpUrl + "/basicset/getall";
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.getall_data = data.Result;
				$scope.college.collegeItems = data.Result[3];
				//console.log("取校学院列表")
				//console.log(data)
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   取科系
	var getMajor = function(collegeid) {
		var temp = [];
		for(var i = 0; i < $scope.getall_data[4].length; i++) {
			if($scope.getall_data[4][i].Collegeid == collegeid) {
				temp.push($scope.getall_data[4][i]);
			}
		}
		$scope.major.majorItems = temp;
	}

	//   取老师
	$scope.queryTeachers = function(collegeid, majorid) {
		Number(collegeid) > 0 ? collegeid = Number(collegeid) : collegeid = 1;
		Number(majorid) > 0 ? majorid = Number(majorid) : majorid = 0;
		var url = config.HttpUrl + "/basicset/queryteachers";
		var data = {
			collegeid: collegeid,
			majorid: majorid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.teacher.teacherItems = data.Result;
				//console.log(data)
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//     select
	$scope.changeSelect = function(title) {
		switch(title) {
			case "college":
				if(!$scope.college.collegeList) {
					$scope.major.majorList = "";
					$scope.major.majorItems = [];
				} else {
					//   科系
					$scope.major.majorList = "";
					getMajor($scope.college.collegeList);
					//  老师
					$scope.queryTeachers($scope.college.collegeList, $scope.major.majorList);
					$scope.teacher.teacherList = "";
				}
				break;
			case "major":
				$scope.queryTeachers($scope.college.collegeList, $scope.major.majorList);
				$scope.teacher.teacherList = "";
				break;
			case "teacher":

				break;
		}
	}

	$scope.ok = function() {
		//console.log(item)
		var temp = {};
		for(var i = 0; i < $scope.teacher.teacherItems.length; i++) {
			if($scope.teacher.teacherItems[i].Usersid == $scope.teacher.teacherList) {
				temp = $scope.teacher.teacherItems[i];
			}
		}
		$modalInstance.close(temp);
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};

	//  run
	var run = function() {
		$scope.getall();
	}
	run();

}]);

/*   弹窗  -  选择班级       */
app.controller("modalGetClassCtrl", ['$scope', 'httpService', '$modalInstance', function($scope, httpService, $modalInstance) {
	//console.log("弹窗-选择班级");

	$scope.getall_data = [];
	//  学院
	$scope.college = {};
	$scope.college.collegeList = "";
	$scope.college.collegeItems = [];
	//   科系
	$scope.major = {};
	$scope.major.majorList = "";
	$scope.major.majorItems = [];
	//   班级
	$scope.class = {};
	$scope.class.classList = "";
	$scope.class.classItems = [];

	//   取校学院列表
	$scope.getall = function() {
		var url = config.HttpUrl + "/basicset/getall";
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.getall_data = data.Result;
				$scope.college.collegeItems = data.Result[3];
				//console.log("取校学院列表")
				//console.log(data)
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   取科系
	var getMajor = function(collegeid) {
		var temp = [];
		for(var i = 0; i < $scope.getall_data[4].length; i++) {
			if($scope.getall_data[4][i].Collegeid == collegeid) {
				temp.push($scope.getall_data[4][i]);
			}
		}
		$scope.major.majorItems = temp;
	}

	//   取班级
	var getClass = function(majorid) {
		var temp = [];
		for(var i = 0; i < $scope.getall_data[5].length; i++) {
			if($scope.getall_data[5][i].Majorid == majorid) {
				temp.push($scope.getall_data[5][i]);
			}
		}
		$scope.class.classItems = temp;
	}

	//     select
	$scope.changeSelect = function(title) {
		switch(title) {
			case "college":
				if(!$scope.college.collegeList) {
					$scope.major.majorList = "";
					$scope.major.majorItems = [];
				} else {
					//   科系
					$scope.major.majorList = "";
					getMajor($scope.college.collegeList);
					//   班级
					$scope.class.classList = "";
					$scope.class.classItems = [];
				}
				break;
			case "major":
				if(!$scope.major.majorList) {
					$scope.class.classList = "";
					$scope.class.classItems = "";
				} else {
					$scope.class.classList = "";
					getClass($scope.major.majorList);
				}
				break;
			case "teacher":

				break;
		}
	}

	$scope.ok = function() {
		//console.log(item)
		var temp = {};
		for(var i = 0; i < $scope.class.classItems.length; i++) {
			if($scope.class.classItems[i].Id == $scope.class.classList) {
				temp = $scope.class.classItems[i];
			}
		}
		$modalInstance.close(temp);
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};

	//  run
	var run = function() {
		$scope.getall();
	}
	run();

}]);

/*   弹窗  -  选择学科       */
app.controller("modalGetScienceCtrl", ['$scope', 'httpService', '$modalInstance', function($scope, httpService, $modalInstance) {
	//console.log("弹窗-选择学科");

	//   一级学科
	$scope.course1 = {};
	$scope.course1.courseList = "";
	$scope.course1.courseItems = [];
	//   二级学科
	$scope.course2 = {};
	$scope.course2.courseList = "";
	$scope.course2.courseItems = [];
	//   三级学科
	$scope.course3 = {};
	$scope.course3.courseList = "";
	$scope.course3.courseItems = [];

	//    
	//    取学科
	var getSubjectClass1 = function(code) {
			var url = config.HttpUrl + "/curriculum/getsubjectclass";
			var data = {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token,
				"SubjectclassCode": code
			};
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				if(data.Rcode == "1000") {
					$scope.course1.courseItems = data.Result;
					//console.log($scope.course1.courseItems)
				} else {
					$scope.tipService.setMessage(data.Reason);
					//alert(data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//    2取学科
	var getSubjectClass2 = function(code) {
			var url = config.HttpUrl + "/curriculum/getsubjectclass";
			var data = {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token,
				"SubjectclassCode": code
			};
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				if(data.Rcode == "1000") {
					$scope.course2.courseItems = data.Result;
					//console.log("取二级学科")
					//console.log($scope.course2.courseItems)
				} else {
					$scope.course2.courseItems = [];
					$scope.tipService.setMessage(data.Reason);
					//alert(data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
		//    3取学科
	var getSubjectClass3 = function(code) {
		var url = config.HttpUrl + "/curriculum/getsubjectclass";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "WEB",
			"Token": config.GetUser().Token,
			"SubjectclassCode": code
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.course3.courseItems = data.Result;
				//console.log("取三级学科")
				//console.log($scope.course3.courseItems)
			} else {
				$scope.course3.courseItems = [];
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//     select
	$scope.changeSelect = function(number) {
		switch(number) {
			case "0":
				if($scope.course1.courseList == "") {
					$scope.course2.courseList = "";
					$scope.course2.courseItems = [];
				} else {
					getSubjectClass2($scope.course1.courseList);
					$scope.course2.courseList = "";
					$scope.course3.courseList = "";
					$scope.course3.courseItems = [];
				}
				break;
			case "1":
				if($scope.course2.courseList == "") {
					$scope.course3.courseList = "";
					$scope.course3.courseItems = [];
				} else {
					getSubjectClass3($scope.course2.courseList);
					$scope.course3.courseList = "";
				}
				break;
			case "2":
				break;
		}
	}

	$scope.ok = function() {
		var temp = {};
		var Subjectcode = "";
		var Subjectname = "";
		var SubjectnameTree = "";
		//  1
		if($scope.course1.courseList != "") {
			for(var i = 0; i < $scope.course1.courseItems.length; i++) {
				if($scope.course1.courseItems[i].Subjectcode == $scope.course1.courseList) {
					SubjectnameTree += $scope.course1.courseItems[i].Subjectname;
					Subjectname = $scope.course1.courseItems[i].Subjectname;
				}
			}
			Subjectcode = $scope.course1.courseList;
		}
		//  2
		if($scope.course2.courseList != "") {
			for(var i = 0; i < $scope.course2.courseItems.length; i++) {
				if($scope.course2.courseItems[i].Subjectcode == $scope.course2.courseList) {
					SubjectnameTree += " - " + $scope.course2.courseItems[i].Subjectname;
					Subjectname = $scope.course2.courseItems[i].Subjectname;
				}
			}
			Subjectcode = $scope.course2.courseList;
		}
		//  3
		if($scope.course3.courseList != "") {
			for(var i = 0; i < $scope.course3.courseItems.length; i++) {
				if($scope.course3.courseItems[i].Subjectcode == $scope.course3.courseList) {
					SubjectnameTree += " - " + $scope.course3.courseItems[i].Subjectname;
					Subjectname = $scope.course3.courseItems[i].Subjectname;
				}
			}
			Subjectcode = $scope.course3.courseList;
		}
		temp = {
			"Subjectcode": Subjectcode,
			"Subjectname": Subjectname,
			"SubjectnameTree": SubjectnameTree
		}
		$modalInstance.close(temp);
	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};

	//  run
	var run = function() {
		getSubjectClass1("");
	}
	run();

}]);

/*   弹窗  -  添加故障         */
app.controller("modalFaultAddCtrl", ['$scope', 'httpService', '$modalInstance', '$modal', 'items','TipService', function($scope, httpService, $modalInstance, $modal, items,TipService) {
	//console.log("弹窗-添加故障");
	
	$scope.tipService = TipService;
	
	//   items==[obj,str]
	$scope.items = items;

		//
	$scope.form = {
		//  添加、查看
		"add": {}
	}
	//   标题
	$scope.title = {
		"add":"故障添加",
		"details":"故障查看",
		"edit":"故障修改"
	}
	//   form
	$scope.form.add = {
		//    标题
		"title":"",
		//   故障id：字符
		"Id": "",
		//   故障设备
		"Device": {},
		//   故障设备
		"DeviceId": "",
		//   故障设备名称
		"DeviceName": "",
		//   教室所有设备Items
		"DeviceItems": [],
		//   设备位置
		"DeviceSite": "",
		//   设备位置 ID
		"DeviceSiteId": "",
		//   故障现象
		"FaultSummary": "",
		//   故障现象 所有
		"FaultSummaryItems": "",
		//   故障描述
		"FaultDescription": "",
		//   故障发生时间
		"HappenTime": "",
		//   设备是否可用  0/1(0-不可使用 1-可以使用)
		"IsCanUse": "0",
		//   申报人 id
		"InputUserId": "",
		//   申报人 名称
		"InputUserName": "",
		//   申报时间
		"InputTime": "",
		//   提交时间（提交故障时间）
		"SubmitTime": "",
		//   故障状态 字符（0-草稿 1-待受理 2-维修中 3-已维修）
		"Status": "0",

		//////故障受理//////
		//   指定维修人
		"AcceptanceRepairPerson": "",
		//   维修人电话
		"AcceptanceRepairPersonTel": "",
		//   受理人ID
		"AcceptanceUserId": null,
		//   受理人姓名
		"AcceptanceUserName": "",
		//   受理人姓名  HTML 显示
		"AcceptanceUserNameHTML": "",
		//   受理时间
		"AcceptanceTime": "",
		
		//////////维修登记//////////
		//   维修人
		"RepairPerson":"",
		//   维修完成时间
		"RepairFinishTime":"",
		//   维修描述
		"RepairDescription":"",
		//   设备是否可用：字符（0不可用，1可以使用）
		"RepairIsCanUse":"0",
		//   维修结果：字符（1-未修复 2-已修复）
		"RepairResult":"1",
		//   维修登记人id：整型
		"RepairInputUserId":null,
		//   维修登记人姓名：字符
		"RepairInputUserName":"",
		//   维修登记人姓名：字符 HTML
		"RepairInputUserNameHTML":"",
		//   维修登记时间：字符
		"RepairInputTime":"",
		//   故障类型
		"RepairFaultType":[],
		//   故障类型 所有
		"RepairFaultTypeItems":[]
	}
	
	$scope.form.add.Device = {
		"DeviceId":"4",
		"DeviceItems":[]
	}
	
	//    有传教室信息，压入对象
	if($scope.items.length == 2){
		$scope.form.add.DeviceId = $scope.items[0].DeviceId;
		$scope.form.add.DeviceSite = $scope.items[0].DeviceSite;
	}

	//   故障记录查询
	$scope.faultItems = {};
	
	//  状态 添加 修改 状态  显示隐藏
	$scope.showSL = false;
	$scope.showDJ = false;
	
	//   add
	$scope.addStatus = false;
	
	//   tab
	$scope.modal_tab = 1;
	$scope.hover = false;
	
	//   查看
	$scope.fault1 = [true,false];
	$scope.fault2 = [true,false];
	$scope.fault3 = [true,false];
	//	验证手机
	$scope.iphone = /^1[34578]{1}[0-9]{9}$/;

	//   生成字符串GUID
	function getGUIDs() {
		var GUID = "";
		for(var i = 1; i <= 32; i++) {
			var n = Math.floor(Math.random() * 16.0).toString(16);
			GUID += n;
			if((i == 8) || (i == 12) || (i == 16) || (i == 20))
				//GUID += "-";
				GUID += "";
		}
		GUID += "";
		return GUID;
	}

	//    打开弹窗  选择教室
	$scope.modalOpenClassroom = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_school.html',
			controller: 'modalGetClassRoomCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			if(!selectedItem) {
				$scope.form.add.DeviceSite = "";
				$scope.form.add.DeviceSite = "";
			} else {
				$scope.form.add.DeviceSite = selectedItem.add;
				$scope.form.add.DeviceSiteId = selectedItem.addId;
				//   教室设备列表
				if(selectedItem.addCode == "classroom") {
					getClassroomDevice($scope.form.add.DeviceSiteId);
				} else {
					alert("请选择具体教室！");
				}
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//   故障发生时间
	$scope.showDate = function() {
			jeDate({
				dateCell: "#happentime",
				format: "YYYY-MM-DD hh:mm:ss",
				isTime: true,
				minDate: "2015-12-31 00:00:00",
				isinitVal: false,
				choosefun: function(elem, val) {
					$scope.form.add.HappenTime = val;
				},
				okfun: function(elem, val) {
					$scope.form.add.HappenTime = val;
				},
				clearfun:function(elem, val) {
					$scope.form.add.HappenTime = "";
				}
			});
		}
	//   申报时间
	$scope.showDate2 = function() {
		jeDate({
			dateCell: "#inputtime",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.form.add.InputTime = val;
			},
			okfun: function(elem, val) {
				$scope.form.add.InputTime = val;
			},
			clearfun:function(elem, val) {
				$scope.form.add.InputTime = "";
			}
		});
	}
	//   维修完成时间
	$scope.showDate3 = function() {
		jeDate({
			dateCell: "#repairfinishtime",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.form.add.RepairFinishTime = val;
			},
			okfun: function(elem, val) {
				$scope.form.add.RepairFinishTime = val;
			}
		});
	}

	//   故障记录查询
	var getFault = function(id) {
		if(!id) return false;
		var url = config.HttpUrl + "/device/getFault";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				Id: id
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.faultItems = data.Result.Data;
				//
				$scope.form.add = $.extend({}, $scope.form.add, $scope.faultItems);
				//console.log($scope.faultItems)
				
				//  设置 维修登记 修改时 默认项
				switch(items[1]) {
					case "add":
						//
					break;
					case "details":
						//
					break;
					case "edit":
						//   设备是否可用：字符（0不可用，1可以使用）
						if($scope.form.add.RepairIsCanUse == '')$scope.form.add.RepairIsCanUse = "0";
						//   维修结果：字符（1-未修复 2-已修复）
						if($scope.form.add.RepairResult == '')$scope.form.add.RepairResult = "1";
					break;
					case "delete":
					break;
				}
				
				//    故障加入
				if($scope.form.add.DeviceItems.length == 0){
					$scope.form.add.DeviceItems = [{"DeviceId":$scope.form.add.DeviceId,"DeviceName":$scope.form.add.DeviceName}];
				}
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    获取教室内所有设备
	var getClassroomDevice = function(classroomid) {
		if(!classroomid) return false;
		var url = config.HttpUrl + "/device/getClassroomDevice";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				ClassroomId: classroomid
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.form.add.DeviceItems = data.Result.Data;
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("获取教室内所有设备")
			//console.log(data)
			//console.log($scope.form.add)
		}, function(reason) {}, function(update) {});
	}

	//   故障登记
	//   ot:操作类型：字符（暂存时传"save",提交时传"submit" 
	var registerFault = function(ot) {
		if(!ot) return false;
		var url = config.HttpUrl + "/device/registerFault";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				//   故障id：字符，不能为空
				"Id": $scope.form.add.Id,
				//   设备id：字符，不能为空
				"DeviceId": $scope.form.add.DeviceId,
				//    故障现象：字符，不能为空
				"FaultSummary": $scope.form.add.FaultSummary,
				//   故障描述：字符，可以为空
				"FaultDescription": $scope.form.add.FaultDescription,
				//   发生时间：字符，不能为空，格式为：yyyy-MM-dd HH:mm:ss
				"HappenTime": $scope.form.add.HappenTime,
				//    是否可用：字符，不能为空（"0"/"1"）
				"IsCanUse": $scope.form.add.IsCanUse,
				//    操作类型：字符（暂存时传"save",提交时传"submit" 
				OT: ot
			}

		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.tipService.setMessage("故障登记成功。");
				//alert("故障登记成功。");
				//
				$modalInstance.close(true);
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   故障受理
	var acceptanceFault = function() {
		var url = config.HttpUrl + "/device/acceptanceFault";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				//故障id：字符，不能为空
				Id: $scope.form.add.Id,
				//维修人：字符，不能为空
				RepairPerson: $scope.form.add.AcceptanceRepairPerson,
				//维修人电话：字符，不能为空
				RepairPersonTel: $scope.form.add.AcceptanceRepairPersonTel
			}

		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.tipService.setMessage("故障受理成功。");
				//alert("故障受理成功。");
				$modalInstance.close(true);
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    维修登记
	var registerRepair = function(ot) {
		var url = config.HttpUrl + "/device/registerRepair";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				//故障id：字符，不能为空
				Id: $scope.form.add.Id, 
				//维修人：字符，不能为空
				RepairPerson: $scope.form.add.RepairPerson, 
				//维修完成时间：字符，不能为空，格式为：yyyy-MM-dd HH:mm:ss
				RepairFinishTime: $scope.form.add.RepairFinishTime, 
				//维修描述：字符，可以为空
				RepairDescription: $scope.form.add.RepairDescription, 
				//维修后设备是否可用：字符，不能为空（"0"/"1"）
				RepairIsCanUse: $scope.form.add.RepairIsCanUse, 
				//维修结果：字符，不能为空（1-未修复 2-已修复）
				RepairResult: $scope.form.add.RepairResult, 
				//故障类型：数组          
				FaultType: $scope.form.add.RepairFaultType,
				//操作类型：字符（暂存时传"save",提交时传"submit"
				OT: ot 
			}

		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.tipService.setMessage("维修登记成功。");
				//alert("维修登记成功。");
				$modalInstance.close(true);
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	
	//   获取设备对应型号的所有故障分类
	var getDeviceAllFaultType = function(deviceid){
		if(!deviceid)return false;
		var url = config.HttpUrl + "/device/getDeviceAllFaultType";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				DeviceId:deviceid
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.form.add.RepairFaultTypeItems = data.Result.Data;
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("所有故障分类")
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	//  取设备故障现象词条
	var getDevicFaultWord = function(deviceid){
		if(!deviceid)return false;
		var url = config.HttpUrl + "/device/getDevicFaultWord";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				DeviceId:deviceid
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.form.add.FaultSummaryItems = data.Result.Data;
			} else {
				$scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("取设备故障词条")
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	//    点选增加现象
	$scope.FaultSummaryClick = function(item){
		var temp;
		if($scope.form.add.FaultSummary == ""){temp = ""}else{temp = "，"};
		$scope.form.add.FaultSummary =  $scope.form.add.FaultSummary + temp + item.Name;
		$scope.hover = false;
	}
	$scope.ToggleClick = function () {
		$scope.hover = !$scope.hover;
	}
	$scope.Mouseleaves = function () {
		$scope.hover = false;
	}

	//   get date
	var getNowFormatDate = function() {
		var date = new Date();
		var seperator1 = "-";
		var seperator2 = ":";
		var month = date.getMonth() + 1;
		var strDate = date.getDate();
		if(month >= 1 && month <= 9) {
			month = "0" + month;
		}
		if(strDate >= 0 && strDate <= 9) {
			strDate = "0" + strDate;
		}
		if(date.getHours() >= 0 && date.getHours() <= 9) {
			var hh = "0" + date.getHours();
		} else {
			var hh = date.getHours();
		}
		if(date.getMinutes() >= 0 && date.getMinutes() <= 9) {
			var mm = "0" + date.getMinutes();
		} else {
			var mm = date.getMinutes();
		}
		if(date.getSeconds() >= 0 && date.getSeconds() <= 9) {
			var ss = "0" + date.getSeconds();
		} else {
			var ss = date.getSeconds();
		}
		var currentdate = date.getFullYear() + seperator1 + month + seperator1 + strDate + " " + hh + seperator2 + mm + seperator2 + ss;
		return currentdate;
	};
	
	//   故障分类标已选中
	var faultFn = function(item,items){
		if(items){
			for(var i = 0; i < items.length; i++){
				items[i].show = true;
				for(var b = 0; b < item.length; b++){
					if(items[i].FaultTypeId == item[b].FaultTypeId){
						items[i].show = false;
						break;
					}
				}
			}
		}
		//   选 中传入RepairFaultType
		var temp = [];
		for(var i = 0; i < items.length; i++){
			if(items[i].show == false){
				temp.push(items[i]);
			}
		}
		$scope.form.add.RepairFaultType = temp;
	}
	
	//   故障分类选中与删除
	$scope.faultTypeClick = function(item,str){
		switch(str){
			case "item":
				//   删除
				for(var i = 0; i < $scope.form.add.RepairFaultTypeItems.length; i++){
					if($scope.form.add.RepairFaultTypeItems[i].FaultTypeId == item.FaultTypeId){
						$scope.form.add.RepairFaultTypeItems[i].show = true;
					}
				}
				var temp = [];
				for(var i = 0; i < $scope.form.add.RepairFaultTypeItems.length; i++){
					if($scope.form.add.RepairFaultTypeItems[i].show == false){
						temp.push($scope.form.add.RepairFaultTypeItems[i]);
					}
				}
				$scope.form.add.RepairFaultType = temp;
				$scope.hover = !$scope.hover;
			break;
			case "items":
				//   添加到选中
				for(var i = 0; i < $scope.form.add.RepairFaultTypeItems.length; i++){
					if($scope.form.add.RepairFaultTypeItems[i].FaultTypeId == item.FaultTypeId){
						$scope.form.add.RepairFaultTypeItems[i].show = false;
					}
				}
				var temp = [];
				for(var i = 0; i < $scope.form.add.RepairFaultTypeItems.length; i++){
					if($scope.form.add.RepairFaultTypeItems[i].show == false){
						temp.push($scope.form.add.RepairFaultTypeItems[i]);
					}
				}
				$scope.form.add.RepairFaultType = temp;
				$scope.hover = !$scope.hover;
			break;
		}
	}
	

	$scope.ok = function(str,number) {
		switch(number) {
			//    故障登记
			case "11":
				if($scope.form.add.Id == "") {
					alert("故障ID不能为空!");
					return false;
				}
				if($scope.form.add.DeviceId == "") {
					alert("故障设备ID不能为空，请选择设备!");
					return false;
				}
				if($scope.form.add.FaultSummary == "") {
					alert("故障现象不能为空，请输入故障现象!");
					return false;
				}
				if($scope.form.add.DeviceSite == "") {
					alert("设备位置不能为空，请选择设备位置!");
					return false;
				}
				if($scope.form.add.HappenTime == "") {
					alert("故障发生时间不能为空，请选择故障发生时间!");
					return false;
				}
				if($scope.form.add.IsCanUse == "") {
					alert("请选择故障设备是否可用!");
					return false;
				}
				//   故障登记
				registerFault(str);
			break;
			//   故障提交  
			case "12":
				if($scope.form.add.Id == "") {
					alert("故障ID不能为空!");
					return false;
				}
				if($scope.form.add.DeviceId == "") {
					alert("故障设备ID不能为空，请选择设备!");
					return false;
				}
				if($scope.form.add.FaultSummary == "") {
					alert("故障现象不能为空，请输入故障现象!");
					return false;
				}
				if($scope.form.add.DeviceSite == "") {
					alert("设备位置不能为空，请选择设备位置!");
					return false;
				}
				if($scope.form.add.HappenTime == "") {
					alert("故障发生时间不能为空，请选择故障发生时间!");
					return false;
				}
				if($scope.form.add.IsCanUse == "") {
					alert("请选择故障设备是否可用!");
					return false;
				}

				//   故障登记
				registerFault(str);
			break;
			//   故障受理
			case "21":
				//console.log($scope.form.add)
				if($scope.form.add.Id == "") {
					alert("故障ID不能为空!");
					return false;
				}
				if($scope.form.add.RepairPerson == "") {
					alert("维修人不能为空!");
					return false;
				}
				if(!($scope.iphone.test($scope.form.add.AcceptanceRepairPersonTel))) {
					alert("维修人电话格式不正确!");
					return false;
				}
				//   故障受理
				acceptanceFault(str);
			break;
			//   维修登记
			case "31":
				//console.log($scope.form.add)
				if($scope.form.add.Id == "") {
					alert("故障ID不能为空!");
					return false;
				}
				if($scope.form.add.RepairPerson == "") {
					alert("维修人不能为空!");
					return false;
				}
				if($scope.form.add.RepairFinishTime == "") {
					alert("维修完成时间不能为空!");
					return false;
				}
				if($scope.form.add.RepairIsCanUse == "") {
					alert("维修设备是否可用不能为空!");
					return false;
				}
				if($scope.form.add.RepairResult == "") {
					alert("维修结果不能为空!");
					return false;
				}
				/*if($scope.form.add.RepairFaultType.length == "") {
					alert("维修故障类型不能为空!");
					return false;
				}*/
				//   
				registerRepair(str);
			break;
			//   维修提交
			case "32":
				//console.log($scope.form.add)
				if($scope.form.add.Id == "") {
					alert("故障ID不能为空!");
					return false;
				}
				if($scope.form.add.RepairPerson == "") {
					alert("维修人不能为空!");
					return false;
				}
				if($scope.form.add.RepairFinishTime == "") {
					alert("维修完成时间不能为空!");
					return false;
				}
				if($scope.form.add.RepairIsCanUse == "") {
					alert("设备是否可用不能为空!");
					return false;
				}
				if($scope.form.add.RepairResult == "") {
					alert("维修结果不能为空!");
					return false;
				}
				/*if($scope.form.add.RepairFaultType.length == "") {
					alert("维修故障类型不能为空!");
					return false;
				}*/
				//   维修完成
				registerRepair(str);
			break;
		}

	};
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};
	
	///////////////////////  添加、修改、查看    状态  //////////////////////////////
	//   items 上级页面传过来状态值
	//  run
	var add = function() {
		//   add
		//   标题
		$scope.form.add.title = $scope.title.add;
		//   隐藏故障受理，维修登记
		$scope.showSL = false;		
		$scope.showDJ = false;
		
		//   故障申报
		$scope.fault1 = [true,false];
		$scope.fault2 = [false,true];
		$scope.fault3 = [false,true];
		
		if(!$scope.items[0]){
			
		}else{
			//     位置
			$scope.form.add.DeviceSite = $scope.items[0].DeviceSite;
			//    设备ID
			$scope.form.add.DeviceId = $scope.items[0].deviceId;
			//  
			getClassroomDevice($scope.items[0].classroomId);
			//    取设备故障词条
			getDevicFaultWord($scope.form.add.DeviceId);
		}
		
		
		//   ID
		$scope.form.add.Id = getGUIDs();
		//   故障发生时间 当前时间
		$scope.form.add.HappenTime = getNowFormatDate();
		//   申报时间
		$scope.form.add.InputTime = getNowFormatDate();
		//   申报人
		$scope.form.add.InputUserName = config.GetUser().Nickname.toString();
		//   申报人ID
		$scope.form.add.InputUserId = config.GetUser().Usersid;

	}
		//add();

	var edit = function() {
		//   标题
		$scope.form.add.title = $scope.title.edit;
		//   修改add状态
		$scope.addStatus = true;
		//   草稿
		if(items[0].Status == 0){
			$scope.modal_tab = 1;
			$scope.fault1 = [true,false];
			$scope.fault2 = [false,true];
			$scope.fault3 = [false,true];
		}
		//   待受理
		if(items[0].Status == 1){
			$scope.modal_tab = 2;
			//$scope.showDJ = true;
			//   故障申报
			$scope.fault1 = [false,true];
			$scope.fault2 = [true,false];
			$scope.fault3 = [false,true];
			
			$scope.form.add.AcceptanceUserNameHTML = config.GetUser().Nickname;
		}
		//   维修中
		if(items[0].Status == 2){
			$scope.modal_tab = 3;
			$scope.fault1 = [false,true];
			$scope.fault2 = [false,true];
			$scope.fault3 = [true,false];
		}
		//   已维修
		if(items[0].Status == 3){
			$scope.fault1 = [false,true];
			$scope.fault2 = [false,true];
			$scope.fault3 = [false,true];
		}
		
		//   位置设备
		if(!$scope.items[0]){
			//getClassroomDevice($scope.items[0].ClassroomId);
		}else{
			//     位置
			$scope.form.add.DeviceSite = $scope.items[0].DeviceSite;
			//    设备ID
			$scope.form.add.DeviceId = $scope.items[0].deviceId;
			//
			$scope.form.add.Id = $scope.items[0].Id;
			//  
			getClassroomDevice($scope.items[0].classroomId);
		}
		
		getFault($scope.form.add.Id);
		//   维修登记
		$scope.form.add.RepairInputUserNameHTML = config.GetUser().Nickname;
	}
	
	var details = function(){
		//   标题
		$scope.form.add.title = $scope.title.details;
		//   查看状态
		$scope.details = true;
		getFault(items[0].Id);
		//   故障申报
		$scope.fault1 = [false,true];
		//   故障受理
		$scope.fault2 = [false,true];
		//   
		$scope.fault3 = [false,true];
		
		//   草稿
		if(items[0].Status == 0){
			$scope.modal_tab = 1;
		}
		//   待受理
		if(items[0].Status == 1){
			$scope.modal_tab = 2;
		}
		//   维修中
		if(items[0].Status == 2){
			$scope.modal_tab = 3;
		}
		//   已维修
		if(items[0].Status == 3){
			
		}
	}
	
	//   指定维修人发生变化
	$scope.$watch('form.add.AcceptanceRepairPerson',function(newValue,oldValue, scope){
		//   维修人
		if(!$scope.form.add.RepairPerson){
			$scope.form.add.RepairPerson = newValue;
		}
	});
	
	//   故障分类
	$scope.$watch('form.add.RepairFaultTypeItems',function(newValue,oldValue, scope){
		//   故障分类选中预处理
		faultFn($scope.form.add.RepairFaultType,newValue);
	});
	
	
	//   设备ID发生变化
	$scope.$watch('form.add.DeviceId',function(newValue,oldValue, scope){
		//   故障分类
		getDeviceAllFaultType($scope.form.add.DeviceId);
		//    取设备故障词条
		getDevicFaultWord($scope.form.add.DeviceId);
	});
	
	//   故障登记人发生变化
	$scope.$watch('form.add.AcceptanceUserName',function(newValue,oldValue, scope){
		//   
		if(newValue != oldValue){
			$scope.form.add.AcceptanceUserNameHTML = $scope.form.add.AcceptanceUserName;
		}
	});
	//   维修登记人发生变化
	$scope.$watch('form.add.RepairInputUserName',function(newValue,oldValue, scope){
		//   
		if(!!$scope.form.add.RepairInputUserName){
			$scope.form.add.RepairInputUserNameHTML = $scope.form.add.RepairInputUserName;
		}
	});
	
	
	///////////////////////  添加、修改、查看  //////////////////////////////

	var run = function() {
		switch(items[1]) {
			case "add":
				add();
			break;
			case "details":
				details();
			break;
			case "edit":
				edit();
			break;
			case "delete":
			break;
		}
	}
	run();

}]);

/*  ------------  弹窗 End  --------------  */