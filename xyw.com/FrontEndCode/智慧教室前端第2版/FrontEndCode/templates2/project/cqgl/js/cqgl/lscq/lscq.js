'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 历史出勤
 */

/*   出勤统计-历史出勤      */
app.controller('cqglLscqContr', ['$scope', 'httpService', '$filter', 'toaster', function($scope, httpService, $filter, toaster) {
	console.log("历史出勤")

	$scope.form = {
		//    开始时间
		"begin_date": "",
		//    结束时间
		"end_date": "",
		//    学院
		"college_ids": "",
		//    专业
		"major_ids": "",
		//    老师
		"teacher_ids": "",
		//    关键词
		"key_word": ""
	}

	$scope.getall_data = [];
	//  学院
	$scope.college = {};
	$scope.college.collegeList = "";
	$scope.college.collegeItem = "";
	$scope.college.collegeItems = [];
	//   科系
	$scope.major = {};
	$scope.major.majorList = "";
	$scope.major.majorItem = "";
	$scope.major.majorItems = [];
	//   老师
	$scope.teacher = {};
	$scope.teacher.teacherList = "";
	$scope.teacher.teacherItem = "";
	$scope.teacher.teacherItems = [];

	//
	$scope.changeTeacherItem = function(item) {
		//
		$scope.changeSelect(item, 'teacher');
	}

	//
	$scope.changeMajorItem = function(item) {
		//
		$scope.changeSelect(item, 'major');
	}

	//
	$scope.changeCollegeItem = function(item) {
		//
		$scope.changeSelect(item, 'college');
	}

	/*  -------------------- 分页、页码  -----------------------  */
	$scope.backPage = {
		"PageCount": 0,
		"PageIndex": 1,
		"PageSize": 15,
		"RecordCount": 0
	};
	/*----------------
	//    分页对象添加页码
	//    return  obj  分页对象
	//    pagedata:obj  分页对象
	//    maxpagenumber:int  显示页码数默认5个页码
	------------------*/
	var pageFn = function(pagedata, maxpagenumber) {
		if(pagedata.length < 1) return null;
		//   缺省时分5页
		Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
		var nub = [];
		var mid = Math.ceil(maxpagenumber / 2);
		if(pagedata.PageCount > maxpagenumber) {
			//  起始页
			var Snumber = 1;
			if((pagedata.PageIndex - mid) < 1) {
				Snumber = 1
			} else if((pagedata.PageIndex + mid) > pagedata.PageCount) {
				Snumber = pagedata.PageCount - maxpagenumber + 1;
			} else {
				Snumber = pagedata.PageIndex - (mid - 1)
			}
			for(var i = 0; i < maxpagenumber; i++) {
				nub.push(Snumber + i);
			}
		} else {
			for(var i = 0; i < pagedata.PageCount; i++) {
				nub.push(i + 1);
			}
		}
		pagedata.Number = nub;
		return pagedata;
	}

	//  翻页
	$scope.pageClick = function(pageindex) {
			if(!(Number(pageindex) > 0)) return false;
			if(pageindex > 0 && pageindex <= $scope.backPage.PageCount) {
				$scope.backPage.PageIndex = pageindex;
				getCurriculums();
			}
		}
		/*  -------------------- 分页、页码  -----------------------  */
		//   开始时间
	$scope.getBeginDate = function() {
			jeDate({
				dateCell: "#lscq_begin",
				format: "YYYY-MM-DD hh:mm:ss",
				isTime: true,
				minDate: "2015-12-31 00:00:00",
				isinitVal: false,
				choosefun: function(elem, val) {
					$scope.form.begin_date = val;
					//
					$scope.searchPost();
				},
				okfun: function(elem, val) {
					$scope.form.begin_date = val;
					//
					$scope.searchPost();
				},
				clearfun: function(elem, val) {
					$scope.form.begin_date = "";
				}
			});
		}
		//   结束时间
	$scope.getEndDate = function() {
		jeDate({
			dateCell: "#lscq_end",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			maxDate: jeDate.now(0),
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.form.end_date = val;
				//
				$scope.searchPost();
			},
			okfun: function(elem, val) {
				$scope.form.end_date = val;
				//
				$scope.searchPost();
			},
			clearfun: function(elem, val) {
				$scope.form.end_date = "";
			}
		});
	}

	//   取校学院列表
	var getall = function() {
		var url = config.HttpUrl + "/basicset/getall";
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.getall_data = data.Result;
				$scope.college.collegeItems = data.Result[3];
				console.log("取校学院列表");
				console.log(data);
			} else {
				toaster.pop('warning', data.Reason);
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
	var queryTeachers = function(collegeid, majorid) {
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
				console.log(data)
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//     select
	$scope.changeSelect = function(item, title) {
		switch(title) {
			case "college":
				if(!$scope.college.collegeList) {
					$scope.major.majorList = "";
					$scope.major.majorItems = [];
					$scope.college.collegeList = item.Id;
					getMajor($scope.college.collegeList);
				} else {
					//   科系
					$scope.major.majorList = "";
				}

				break;
			case "major":
				$scope.college.collegeList = item.Collegeid;
				$scope.major.majorList = item.Majorid;
				queryTeachers($scope.college.collegeList, $scope.major.majorList);
				$scope.teacher.teacherList = "";
				break;
			case "teacher":
				$scope.teacher.teacherList = item.Usersid;
				$scope.college.collegeList = item.Collegeid;
				$scope.major.majorList = item.Majorid;
				queryTeachers($scope.college.collegeList, $scope.major.majorList);
				break;
		}

		//
		$scope.searchPost();
	}

	//  用户查课表  当天内
	var getCurriculums = function() {
		var myDate = new Date();
		var temp_end = "";
		new Date($scope.form.end_date) > myDate ? temp_end = $filter('date')(myDate, 'yyyy-MM-dd HH:mm:ss') : temp_end = $scope.form.end_date;

		var url = config.HttpUrl + "/action/gethistoryattendance";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Begindate": $scope.form.begin_date,
			"Enddate": temp_end,
			"State": -1,
			"Collegeids": $scope.form.college_ids,
			"Majorids": $scope.form.major_ids,
			"Teacherids": $scope.form.teacher_ids,
			"Searhtxt": $scope.form.key_word,
			"PageSize": Number($scope.backPage.PageSize),
			"PageIndex": Number($scope.backPage.PageIndex)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.schedule = data.Result.PageData;
				console.log($scope.schedule);
				//   出勤率
				for(var a in $scope.schedule) {
					$scope.schedule[a].Toclassrate = Math.round($scope.schedule[a].Toclassrate * 1000) / 10.0 + "%";
				}
				var objPage = {
					PageCount: 0,
					PageIndex: data.Result.PageIndex,
					PageSize: data.Result.PageSize,
					RecordCount: data.Result.PageCount
				};
				if((objPage.RecordCount % objPage.PageSize) == 0) {
					objPage.PageCount = (objPage.RecordCount / objPage.PageSize);
				} else {
					objPage.PageCount = parseInt((objPage.RecordCount / objPage.PageSize)) + 1;
				}
				//   分页
				//$scope.backPage.PageIndex = $scope.schedule;
				$scope.backPage = pageFn(objPage, 5);
			} else {
				$scope.backPage.PageCount = 0;
				$scope.schedule = "";
        toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   学院
	$scope.$watch('college.collegeList', function(newValue, oldValue, scope) {
		$scope.form.college_ids = newValue;
		console.log($scope.form)
	});
	//   科系
	$scope.$watch('major.majorList', function(newValue, oldValue, scope) {
		$scope.form.major_ids = newValue;
		console.log($scope.form)
	});
	//   老师
	$scope.$watch('teacher.teacherList', function(newValue, oldValue, scope) {
		$scope.form.teacher_ids = newValue;
		console.log($scope.form)
	});

	//   查询
	$scope.searchPost = function() {
		$scope.backPage.PageIndex = 1;
		getCurriculums();
	}

	//   回车查询
	$scope.sbgzKeyup = function(e) {
		var keycode = window.event ? e.keyCode : e.which;
		if(keycode == 13) {
			$scope.searchPost();
		}
	}

	//   run
	var run = function() {
		//   当天
		var myDate = new Date();

		//   监听服务器时间 变化
		//   只执行一次
		var one_s = true;
		$scope.$watch('app.serverTime', function(newValue, oldValue, scope) {
			if(newValue && one_s){
				$scope.form.begin_date = $filter('date')($scope.app.serverTime * 1000, 'yyyy-MM-dd') + " 00:00:00";
				$scope.form.end_date = $filter('date')($scope.app.serverTime * 1000, 'yyyy-MM-dd HH:mm:ss');
				one_s = false;
			}
		});
		//
		getCurriculums();

		getall();
	}
	run();
}]);

/*   历史出勤-查看     */
app.controller("cqglLscqDetailsContr", ['$scope', '$location', 'httpService','toaster', function($scope, $location, httpService,toaster) {
	console.log("历史出勤-查看")

	$scope.CId = $location.search().CId;
	//   单条课表
	$scope.CurriculumsInfo = {};

	$scope.form = {
		//------------------//
		//   上课日期
		"day": "",
		//    上课时间
		"time": "",
		//    教室
		"room": "",
		//   上课老师
		"teacher": "",
		//    课程
		"lesson": "",
		//---------------------//
		//    上课班级
		"class_room": "",
		//    应到人数
		"due_people": "",
		//    实到人数
		"actual_people": "",
		//    缺勤人数
		"lack_people": "",
		//    出勤率
		"cq_people": ""
	}

	//  查单条课表
	var getCurriculumsInfo = function(cid) {
		if(Number(cid) > -1) {
			cid = Number(cid)
		} else {
			return false
		}
		var url = config.HttpUrl + "/action/getcurriculumsinfo";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Curriculumclassroomchaptercentreid": cid
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("查单条课表", data)
			if(data.Rcode == "1000") {
				$scope.CurriculumsInfo = data.Result;
				//   对象不为空
				if($scope.CurriculumsInfo != null && !angular.equals({}, $scope.CurriculumsInfo)) {
					//   上课日期
					$scope.form.day = $scope.CurriculumsInfo.Begindate.substr(0, 10);
					//   time
					$scope.form.time = $scope.CurriculumsInfo.Begindate.substr(11, 5) + "-" + $scope.CurriculumsInfo.Enddate.substr(11, 5);
					//   教室
					$scope.form.room = $scope.CurriculumsInfo.Buildingname + "-" + $scope.CurriculumsInfo.Campusname + "-" + $scope.CurriculumsInfo.Buildingname + "-" + $scope.CurriculumsInfo.Classroomsname;
					//   上课老师
					$scope.form.teacher = $scope.CurriculumsInfo.Truename;
					//   课程
					$scope.form.lesson = $scope.CurriculumsInfo.Curriculumname;
					//   上课班级
					$scope.form.class_room = $scope.CurriculumsInfo.Classesname;
					//   应到人数
					$scope.form.due_people = $scope.CurriculumsInfo.Plannumber;
					//   实到人数
					$scope.form.actual_people = $scope.CurriculumsInfo.Actualnumber;
					//   缺勤人数
					$scope.form.lack_people = $scope.CurriculumsInfo.Plannumber - $scope.CurriculumsInfo.Actualnumber;
					//   出勤率
					$scope.form.cq_people = $scope.form.due_people > 0 ? ($scope.form.actual_people / $scope.form.due_people).toFixed(2) + "%" : "0%";
				}
			} else {
        toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    get 点到学生信息
	var getPointtos = function(id) {
		if(Number(id) < 0 && !id) {
			return false
		} else {
			id = Number(id)
		};
		//   不在上课时间 id 为0
		if(id == 0) {
			return false
		}
		var url = config.HttpUrl + "/action/getpointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Curriculumclassroomchaptercentreid": id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.pointtos = data.Result;
				//   实到 与 未到
				if($scope.pointtos.length > 0) {
					var temp = 0;
					for(var a in $scope.pointtos) {
						if($scope.pointtos[a].State == 0) {
							temp += 1;
						}
					}
					//   应到人数
					$scope.form.due_people = $scope.pointtos.length;
					//   实到
					$scope.form.actual_people = $scope.form.due_people - temp;
					//   未到
					$scope.form.lack_people = temp;
					//   出勤率
					$scope.form.cq_people = Math.round($scope.form.actual_people / $scope.form.due_people * 1000) / 10.0 + "%";
				}
			} else {
        toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	var run = function() {
		getCurriculumsInfo($scope.CId);
		getPointtos($scope.CId);
	}
	run();
}]);
