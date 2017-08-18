'use strict';
/**
 * Created by Administrator on 2016/8/17.
 * 中控面板
 */

app.controller("zkmbindexContr", ['$scope', '$rootScope', '$location', 'httpService', '$modal', '$timeout', '$interval', '$state', '$filter','$window', function($scope, $rootScope, $location, httpService, $modal, $timeout, $interval, $state, $filter,$window) {

	//  默认值
	var timer;
	//   默认图片路径
	$scope.deviceimg = config.zkmb_config.deviceImg;

	//   隐藏公共部分
	$rootScope.showcom = false;

	//   退出
	$scope.loginout = function() {
		localStorage.removeItem("LoginUser");
		window.location.href = "/web2/html/login_zkmb.html";
	};

	//   上课
	$scope.attendClass = function() {
		var url = config.HttpUrl + "/action/changeclassstate";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Classroomid": Number(config.zkmb_config.classroomId),
			//   	[0:上课,1:下课]
			"State": 0,
			"Ccccids":config.zkmb_config.Ccccid
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("上课",data)
			if(data.Rcode == "1000"){
				
				//   上课状态
				var temp = "";
				for(var a in data.Result){
					temp += data.Result[a].Ccccid + ",";
				}
				temp.length > 2 ?　temp = temp.substr(0,temp.length - 1) : temp = temp;
				
				$scope.attend.inCcccid = temp;
				
				config.zkmb_config.Ccccid = temp;
				$scope.tipService.setMessage("上课成功！");
			}else{
				//   上课不成功继续执行上课动作开关
				//classroom_one = true;
				//$scope.tipService.setMessage(data.Reason);
			}
			//   有上课动作  可以有执行下课动作
			//classroom_noup = true;
		}, function(reason) {}, function(update) {});
	}
	
		
		//   下课
	$scope.attendClassOut = function() {
		var url = config.HttpUrl + "/action/changeclassstate";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Classroomid": Number(config.zkmb_config.classroomId),
			//   	[0:上课,1:下课]
			"State": 1,
			"Ccccids":$scope.attend.inCcccid
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("下课",data)
			if(data.Rcode == "1000") {
				
				//   下课状态
				var temp = "";
				for(var a in data.Result){
					temp += data.Result[a].Ccccid + ",";
				}
				temp.length > 2 ?　temp = temp.substr(0,temp.length - 1) : temp = temp;
				
				$scope.attend.outCcccid = temp;
				
				//    下课清空
				config.zkmb_config.Ccccid = "";
				
				$scope.tipService.setMessage("下课成功！");
			} else{
				//   下课不成功继续执行上课动作开关
				//classroom_one = true;
				//$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    下课
	$scope.outCourse = function() {
		var modalInstance = $modal.open({
			templateUrl: '../project/zkmb/html/zkmb/modal_exit.html',
			windowClass: "m-zkmb-modal",
			animate:false,
			controller: "exitZkmbContr"
		});
		//$scope.loginout();
	}

	//   取教室
	$scope.zkmb_class_info = [];
	$scope.getClassroomInfo = function() {
		var url = config.HttpUrl + "/basicset/getclassroominfo";
		var data = {
			id: config.zkmb_config.classroomId
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			console.log("取教室教室信息",data)
			if(data.Rcode == 1000){
				$scope.zkmb_class_info = data.Result;
				//    班级ID
				config.zkmb_config.classId = $scope.zkmb_class_info.Classesid;
			}else{
				
			}
			//console.log("取教室",data)
		}, function(reason) {}, function(update) {})
	}
	

	//获取设备状态信息
	$scope.zkmb_fnGetDeviceStatus = function() {
			var url = config.zkmb_config.coapServer + "/device/node/state/room";
			var data = {
				UserID: config.zkmb_config.Uid, 		//用户ID(*字符串)   
				RoomID: config.zkmb_config.classroomId, 	//房间ID(*字符串)    
			    Params: "" 	//参数(字符串)
			};
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				$scope.deviceData = data.Data;
				//console.log($scope.deviceData)
				//fnAutoRefresh();
				//   自动刷新
				tm.fnAutoRefreshfn(tm);
				$.unblockUI();
			}, function(reason) {}, function(update) {})
		}
	

	// 教室设备 一键开启
	$scope.fnOpenDevice = function() {
		tm.fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/on/room";
		var data = {
			"UserID": config.zkmb_config.Uid,	//(*字符串)
			"RoomID": config.zkmb_config.classroomId, //(*字符串)
			"Params": ""	//(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			tm.fnAutoRefresh();
		}, function(reason) {}, function(update) {})
	}

	//   教室设备 一键关闭
	$scope.fnCloseDevice = function() {
		tm.fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/off/room";
		var data = {
			"UserID": config.zkmb_config.Uid,	//(*字符串)
			"RoomID": config.zkmb_config.classroomId, //(*字符串)
			"Params": ""	//(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			tm.fnAutoRefresh();
		}, function(reason) {}, function(update) {})
	}


	
	
	//开始定义定时器
	var tm=$scope.setglobaldata.gettimer("zkmb_index_device");
	if(tm.Key!="zkmb_index_device"){
		tm.Key="zkmb_index_device";
		tm.keyctrl="app.zkmb";
		tm.fnAutoRefresh=function(){
			//console.log("开始调用定时器");
			this.interval = $interval(function() {
				$scope.zkmb_fnGetDeviceStatus();
			},config.zkmb_config.fixRefreshTime);	
		};
		tm.fnStopAutoRefresh=function(){
			//console.log("进入取消方法");
			if(!angular.isUndefined(this.interval)) {
				$interval.cancel(this.interval);
				this.interval = 'undefined';
				//console.log("进入取消成功");
			}
			this.interval=null;
		};
		$scope.setglobaldata.addtimer(tm);
	}
	//结束定义定时器
	
	

	//  打开
	$scope.zkmb_sbkz = function(ite) {
		config.zkmb_config.item = ite;
		$state.go("app.zkmb.sbkz", {
			"DeviceId": ite.DeviceId,
			"page": ite.DevicePage
		});
	}

	/*   //////////////////////////////////////  */
	//   返回
	$scope.goback = function() {
		$state.go('^');
	}

	//  判断不是在中控面板首页
	$scope.isZkmb = true;
	$rootScope.$on('$stateChangeStart',
		function(event, toState, toParams, fromState, fromParams) {
			if(toState.name != "app.zkmb") {
				$scope.isZkmb = false;
			} else {
				$scope.isZkmb = true;
				//  run
				$scope.zkmb_fnGetDeviceStatus();
			}
		});
		
	//    手动刷新处理
	if($state.current.name != "app.zkmb") {
		$scope.isZkmb = false;
	} else {
		$scope.isZkmb = true;
		//  run
		$scope.zkmb_fnGetDeviceStatus();
	}
	
	
	
	

	///////////////////////////////////////////////////////
	//   课表
	$scope.curriculums = [];
    //   查课表日期
	$scope.curriculumsDay = $filter('date')(new Date(), 'yyyy-MM-dd');
    //   当前课程
	$scope.thisCourse = {};
	//   学生出勤
	$scope.pointtos = [];
	//   考勤
	$scope.attendance = {
		//  应到
		all: 0,
		//  实到
		act: 0,
		/*//   百分比
		percentage: 0*/
	};
	//   上下课状态
	$scope.attend = {
		//  上课状态
		inCcccid:"0",
		//  下课状态
		outCcccid:"0"
	}
	//   当前上课班级
	$scope.inClass = [];
	
	//   用户与老师不相同false  相同true
	$scope.isTeacher = false;

	//   换课
	$scope.zkmb_huanke = function(ite) {
		var modalInstance = $modal.open({
			templateUrl: '../project/zkmb/html/zkmb/modal_huanke.html',
			controller: 'zkmbHuankeindexContr',
			resolve: {
				items: function() {
					return $scope.curriculums;
				}
			}
		});
	}

	//   课间休息
	//   classRecessArray[3] == 0;
	var classRecessArray = [true, "课间休息", "继续上课",0];
	$scope.classRecessText = "课间休息";
	$scope.classRecess = function() {
		if(classRecessArray[0]) {
			$scope.classRecessText = classRecessArray[2];
			classRecessArray[0] = false;
			classRecessArray[3] = 2;
			//   下课
			$scope.attendClassOut();
		} else {
			$scope.classRecessText = classRecessArray[1];
			classRecessArray[0] = true;
			classRecessArray[3] = 1;
			//   上课
			$scope.attendClass();
		}
	}

	//  用户查课表  当天内
	var getCurriculums = function() {
		var myDate = new Date();
		var url = config.HttpUrl + "/action/getcurriculums";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Begindate": $scope.curriculumsDay + " 00:00:00",
			"Enddate": $scope.curriculumsDay + " 23:59:59",
			"State": -1,
			"Classroomid":Number(config.zkmb_config.classroomId),
			"PageIndex":-1
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("今天课表",data)
			if(data.Rcode == "1000") {
				$scope.curriculums = data.Result;
				getCourse($scope.curriculums);
				//   定时器
				tm2.fnAutoRefreshfn(tm2);
			} else {
				//$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	

	//   课表查找当前时间课程  验证是否在上课时间
	var getCourse = function(items) {
		if(!angular.isObject(items)) return false;
		var myDate = new Date();
	    //   是否在当天
		if ($filter('date')(myDate, 'yyyy-MM-dd') != $scope.curriculumsDay) {
            //   当前课表不在当天 重新查课表
		    $scope.curriculumsDay = $filter('date')(new Date(), 'yyyy-MM-dd');
		    getCurriculums();
		    //    重新取服务器时间
		    $scope.getServerTime();
		}
	    //  时间戳 秒
		var t = $scope.app.serverTime ? $scope.app.serverTime : Date.parse(myDate) / 1000;
		//    当前上课班级字符串,隔开
		var Curriculumclassroomchaptercentreids="";
		//    当前上课班级数组
		$scope.inClass = [];
		for (var i = 0; i < items.length; i++) {
            //   begindate
		    var b = Date.parse(new Date(items[i].Begindate.replace(/\-/g, "/"))) / 1000;
            //   enddate
		    var e = Date.parse(new Date(items[i].Enddate.replace(/\-/g, "/"))) / 1000;
			if(t >= b && t <= e) {
				//   在上课时间  用户不同
				if(items[i].TeacherId != config.GetUser().Usersid){
					$scope.thisCourse = "";
					$scope.isTeacher = false;
					
					//   传入视频录制参数
					$scope.formVideo.CurriculumName = "";
					$scope.formVideo.TeacherName = "";
					$scope.formVideo.CharpterName = "";
					$scope.formVideo.CurriculumDuration = null;
					$scope.formVideo.CurriculumClassroomChapterID = null;
				}else{
					//  课表上课时间内 当前时间课程 ,当前用户对应课表老师
					$scope.thisCourse = items[i];
					
					console.log("当前上课",items[i]);
					$scope.isTeacher = true;
					$scope.inClass.push($scope.thisCourse);
					Curriculumclassroomchaptercentreids=Curriculumclassroomchaptercentreids+$scope.thisCourse.Curriculumclassroomchaptercentreid+",";
					
					//    传入视频录制参数
					//   begindate 秒
				    var b = Date.parse(new Date($scope.thisCourse.Begindate.replace(/\-/g, "/"))) / 1000;
		            //   enddate 秒
				    var e = Date.parse(new Date($scope.thisCourse.Enddate.replace(/\-/g, "/"))) / 1000;
					$scope.formVideo.CurriculumName = $scope.thisCourse.Curriculumname;
					$scope.formVideo.TeacherName = $scope.thisCourse.Truename;
					$scope.formVideo.CharpterName = $scope.thisCourse.Chaptername;
					$scope.formVideo.CurriculumDuration = e - b;
					$scope.formVideo.CurriculumDurationEnddDate = $scope.thisCourse.Enddate;
					$scope.formVideo.CurriculumClassroomChapterID = $scope.thisCourse.Curriculumclassroomchaptercentreid;
					
				}
				
				//    没点过课间休息 && 是继续上课
				if(classRecessArray[3] == 0 && classRecessArray[0] == false){
					//    改成课间休息
					$scope.classRecessText = classRecessArray[1];
					classRecessArray[0] = true;
					classRecessArray[3] = 1;
				}
			}
		}
		//
		Curriculumclassroomchaptercentreids=Curriculumclassroomchaptercentreids.substring(0,Curriculumclassroomchaptercentreids.length-1);
		if(Curriculumclassroomchaptercentreids!="" && Curriculumclassroomchaptercentreids!=undefined){
			$scope.thisCourse.Curriculumclassroomchaptercentreids=Curriculumclassroomchaptercentreids;
			$scope.getPointtos(Curriculumclassroomchaptercentreids);
		}else{
			//    不在上课时间
			$scope.thisCourse = "";
			//    不在上课时间 不能纠正
			$scope.isTeacher = false;
			//    清空应到实到
			$scope.attendance.all = 0;
			$scope.attendance.act = 0;
			//    不在课时,清教室字符
			$scope.attend.inCcccid = "";
			config.zkmb_config.Ccccid = "";
			//    圆
			var temp = angular.copy($scope.option_pie);
				temp.series[0].data[0].value = 0;
				temp.series[0].data[1].value = 1;
				$scope.option_pie = temp;
				
			//    课间休息 状态 改成 没点过
			classRecessArray[3] = 0;
		}
		//   上课
		if(($scope.attend.inCcccid == "" || $scope.attend.inCcccid == null) && $scope.isTeacher){
			//  上课
			$scope.attendClass();
		}
	}

	//    get 点到
	$scope.getPointtos = function(id) {
		if(!id)return false;
		//console.log(id);
		var myDate = new Date();
		var url = config.HttpUrl + "/action/getpointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			//"Curriculumclassroomchaptercentreid": Number(id)
			"Curriculumclassroomchaptercentreids":id.toString()
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("get 点到",data)
			if(data.Rcode == "1000") {
				$scope.pointtos = data.Result;

				//   应到
				$scope.attendance.all = $scope.pointtos.length;
				
				//   实到
				$scope.attendance.act = 0;
				for(var i = 0; i < $scope.pointtos.length; i++) {
					if($scope.pointtos[i].State == 1) {
						$scope.attendance.act += 1;
					}
				}
				var temp = angular.copy($scope.option_pie);
				temp.series[0].data[0].value = $scope.attendance.act;
				temp.series[0].data[1].value = $scope.attendance.all - $scope.attendance.act;
				$scope.option_pie = temp;
				
			} else {
				$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	

	//    echars 圆
	$scope.option_pie = {
		color:['#f86896','#ffffff'],
		series: [
			{
				name:'上课点到',
				type:'pie',
				radius: ['87%', '100%'],
				label: {
					normal: {
						show: false
					}
				},
				animation:false,
				data:[
					{value:0},
					{value:1}
				]
			}
		]
	};
	


	/*    /////////////////////////////////////    */

	
	

	/*  ///////////////////////////////////////////////  */

	//  录制
	$scope.video = {
		//   录制状态  false 关闭
		record: false,
		//   录制时长
		time: 0,
		//   开始时间
		startTime: "",
		//   结束时间
		stopTime: "",
		//   定时器
		timer: null
	}
	//   视频表单
	$scope.formVideo = {
		//   课程名
		"CurriculumName":"",
		//   教师名称
		"TeacherName":"",
		//   章节名称
		"CharpterName":"",
		//   课程时长(秒)
		"CurriculumDuration":null,
		//   课程结束时间
		"CurriculumDurationEnddDate":"",
		//   教室ID
		"ClassroomID":config.zkmb_config.classroomId,
		//   课程教室章节ID
		"CurriculumClassroomChapterID":null
	}

	//   开始录制时间、时长
	var startRecord = function() {
		var myDate = new Date();
		$scope.video.record = true;
		$scope.video.time = 0;
		$scope.video.startTime = myDate;
		$scope.video.stopTime = "正在录制";
		$scope.video.timer = $interval(function() {
			$scope.video.time += 1;
		}, 1000);
	}

	//   结束录制时间、时长
	var stopRecord = function() {
		var myDate = new Date();
		$scope.video.record = false;
		$scope.video.stopTime = myDate;
		$interval.cancel($scope.video.timer);
	}

	//    开始录制
	$scope.beginvideo = function() {
		//   时长
		var myDate = new Date();
		var begin = $scope.app.serverTime ? $scope.app.serverTime : Date.parse(myDate) / 1000;
		var end = Date.parse(new Date($scope.thisCourse.Enddate.replace(/\-/g, "/"))) / 1000;
		if(end - begin > 0){
			$scope.formVideo.CurriculumDuration = end - begin;
		}else{
			return false;
		}
		//   手动设置5分钟
		$scope.formVideo.CurriculumDuration = 60;
		
		var url2 = config.HttpUrl + "/action/beginvideo";
		var data2 = {
			"UsersID": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "WEB",
			"Token": config.GetUser().Token,
			
			"CurriculumName": $scope.formVideo.CurriculumName.toString(),
			"TeacherName": $scope.formVideo.TeacherName.toString(),
			"CharpterName": $scope.formVideo.CharpterName.toString(),
			"CurriculumDuration": $scope.formVideo.CurriculumDuration,
			"ClassroomID": Number($scope.formVideo.ClassroomID),
			"CurriculumClassroomChapterID":Number($scope.formVideo.CurriculumClassroomChapterID)
		};
		//return false;
		var promise2 = httpService.ajaxPost(url2, data2);
		promise2.then(function(data) {
			console.log("开始录制",data)
			if(data.Rcode == "1000") {
				//   开始时间
				startRecord();
				//   自动停止
				var tempTime = Number(angular.copy($scope.formVideo.CurriculumDuration));
				var tempTimeFn = function(){
					$timeout(function(){
						if(tempTime > 0){
							tempTime -= 1;
							tempTimeFn();
						}else{
							stopRecord();
						}
					},1000);
				}
				tempTimeFn();
			} else {
				$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	//  结束录制
	$scope.endvideo = function() {
		//   时长
		var myDate = new Date();
		var begin = $scope.app.serverTime ? $scope.app.serverTime : Date.parse(myDate) / 1000;
		var end = Date.parse(new Date($scope.thisCourse.Enddate.replace(/\-/g, "/"))) / 1000;
		$scope.formVideo.CurriculumDuration = end - begin;
		
		var url2 = config.HttpUrl + "/action/endvideo";
		var data2 = {
			"UsersID": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "WEB",
			"Token": config.GetUser().Token,
			
			"CurriculumName": $scope.formVideo.CurriculumName.toString(),
			"TeacherName": $scope.formVideo.TeacherName.toString(),
			"CharpterName": $scope.formVideo.CharpterName.toString(),
			"CurriculumDuration": $scope.formVideo.CurriculumDuration,
			"ClassroomID": Number($scope.formVideo.ClassroomID)
		};
		var promise2 = httpService.ajaxPost(url2, data2);
		promise2.then(function(data) {
			console.log("结束录制",data)
			if(data.Rcode == "1000") {
				//   结束时间
				stopRecord();
				//console.log(data.Rcode + "--" + data.Reason);
			} else {
				$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	/*  ///////////////////////////////////////////////  */

	//   后退出中控面板
	$scope.$on("$destroy", function() {
		$rootScope.showcom = true;
	});
	
	//    不同栏目不同标题
	$scope.sbkz_title_show = false;
	//    返回刷新点到  监听路由
	$rootScope.$on('$stateChangeSuccess',
		function(event, toState, toParams, fromState, fromParams) {
			if(toState.name == "app.zkmb"){
				getCourse($scope.curriculums);
			}
			if(toState.name == "app.zkmb.sbkz"){
				$scope.sbkz_title_show = true;
			}else{
				$scope.sbkz_title_show = false;
			}
		}
	);
	//    手动刷新处理
	if($state.current.name == "app.zkmb"){
		getCourse($scope.curriculums);
	}
	if($state.current.name == "app.zkmb.sbkz"){
		$scope.sbkz_title_show = true;
	}else{
		$scope.sbkz_title_show = false;
	}
	
	
	//开始定义定时器
	var tm2 = $scope.setglobaldata.gettimer("zkmb_index_dd");
	if(tm2.Key!="zkmb_index_dd"){
		tm2.Key="zkmb_index_dd";
		tm2.keyctrl="app.zkmb";
		tm2.fnAutoRefresh=function(){
			this.interval = $interval(function() {
				//   查课表
				getCourse($scope.curriculums);
				
				//buildingString();//   格式化字符串
				//getClassroomStatusList(buildingids);//   get 教室
			}, config.zkmb_config.zkmbRefreshTime);	
		};
		tm2.fnStopAutoRefresh=function(){
			console.log("进入取消方法");
			if(!angular.isUndefined(this.interval)) {
				$interval.cancel(this.interval);
				this.interval = 'undefined';
				console.log("进入取消成功");
			}
			this.interval=null;
		};
		$scope.setglobaldata.addtimer(tm2);
	}
	//结束定义定时器
	
	
	/*  ------------  设置  开始   ---------------  */
	
	
	
	
	/**
	 * 中控面板设置教室 打开弹窗
	 */
	$scope.zkmb_sys_option = function(){
		var modalInstance = $modal.open({
			templateUrl: '../project/zkmb/html/zkmb/modal_option.html',
			controller: 'modalGetClassRoomCtrl',
			windowClass: 'm-modal-option',
			animate:false,
			//size: size,
			resolve: {
				items: function() {
					return '';
				}
			}
		});
		
		//   
		modalInstance.result.then(function(item) {
			//
			console.log(item);
			if(item){
				if(('addCode' in item) && item.addCode == 'classroom'){
					$window.localStorage['zkmb_option'] = JSON.stringify(item);
					config.zkmb_config.classroomId = item.addId;
					//   初始化视频
					$scope.formVideo.ClassroomID = config.zkmb_config.classroomId;
					$scope.getClassroomInfo();
				}
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	
	/*  ------------  设置  结束  ---------------  */
	
	
	//   run
	var run = function(){
		//   是否已选择地址教室
		if(!$window.localStorage['zkmb_option']){
			$scope.zkmb_sys_option();
		}else{
			config.zkmb_config.classroomId = JSON.parse($window.localStorage['zkmb_option']).addId;
			//   初始化视频
			$scope.formVideo.ClassroomID = config.zkmb_config.classroomId;
		}
		
		//   取教室信息
		$scope.getClassroomInfo();
		//   取设备信息
		$scope.zkmb_fnGetDeviceStatus();
		
		//   取课表
		getCurriculums();
		//   上课
		//$scope.attendClass();
		//   下课
		//$scope.attendClassOut();
	}
	run();

}]);




/*   中控弹窗    */
app.controller("exitZkmbContr", ['$scope', '$modalInstance','httpService','TipService','$filter', function($scope, $modalInstance,httpService,TipService,$filter) {
	//$scope.zkmb_config = items;
	//   中控弹窗
	$scope.tipService = TipService;
	//   下课 点取消
	$scope.cancel = function() {
		$modalInstance.close();
	}
	$scope.loginout = function() {
		$scope.classOut = true;
		//   下课
		var url = config.HttpUrl + "/action/changeclassstate";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Classroomid": Number(config.zkmb_config.classroomId),
			//   	[0:上课,1:下课]
			"State": 1,
			"Ccccids":config.zkmb_config.Ccccid
		};
		//   超时   1000 毫秒
		var promise = httpService.ajaxPost(url, data,1000);
		promise.then(function(data) {
			$scope.classOut = false;
			if(data.Rcode == "1000") {
				$scope.tipService.setMessage(data.Reason);
			} else{
				//$scope.tipService.setMessage(data.Reason);
			}
			localStorage.removeItem("LoginUser");
			window.location.href = "/web2/html/login_zkmb.html";
		}, function(reason) {}, function(update) {});
	}
}]);
/*   换课    */
app.controller("zkmbHuankeindexContr", ['$scope', 'items', function($scope, items) {
	$scope.huankeItem;
	$scope.huankeItems = items;
}]);

/*  考勤纠正    */
app.controller("zkmbKqjzindexContr", ['$scope', '$location', '$rootScope', 'httpService', '$filter', '$log','$interval', function($scope, $location, $rootScope, httpService, $filter, $log,$interval) {
	console.log("考勤纠正")
	$scope.pointtos = [];

	//    get 点到
	$scope.getPointtos = function(id) {
		if(!id)return false;
		var url = config.HttpUrl + "/action/getpointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			//"Curriculumclassroomchaptercentreid": Number(id)
			"Curriculumclassroomchaptercentreids":id.toString()
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log("考勤纠正-点到",data)
			if(data.Rcode == "1000") {
				$scope.pointtos = data.Result;
				//   默认点到按钮 显示第一个
				for(var i = 0; i < $scope.pointtos.length; i++) {
					$scope.pointtos[i].show = false;
				}
			} else {
				//$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
		
	}
		

	//   点到
	$scope.postPointtos = function(item) {
		if(!item)return false;
		$.blockUI({
			message: '<div class="model_blockui" style="padding: 10px;"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;ִ点到中...</span></div>',timeout:10000
		});
		var url = config.HttpUrl + "/action/updatepointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "WEB",
			"Token": config.GetUser().Token,
			//"Curriculumclassroomchaptercentreid": Number($scope.thisCourse.Curriculumclassroomchaptercentreid),
			"Curriculumclassroomchaptercentreids":$scope.thisCourse.Curriculumclassroomchaptercentreids,
			"Studentsid": Number(item.Usersid),
			"State": 1
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreids);
			} else {
				$scope.tipService.setMessage(data.Reason);
			}
			$.unblockUI();
		}, function(reason) {}, function(update) {});
	}

	//   取消点到
	$scope.postPointtosClose = function(item) {
		$.blockUI({
			message: '<div style="padding: 10px;"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;ִ点到中...</span></div>',timeout:10000
		});
		var url = config.HttpUrl + "/action/updatepointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "WEB",
			"Token": config.GetUser().Token,
//			"Curriculumclassroomchaptercentreid": Number($scope.thisCourse.Curriculumclassroomchaptercentreid),
			"Curriculumclassroomchaptercentreids":$scope.thisCourse.Curriculumclassroomchaptercentreids,
			"Studentsid": Number(item.Usersid),
			"State": 0
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreids);
			} else {
				$scope.tipService.setMessage(data.Reason);
			}
			$.unblockUI();
		}, function(reason) {}, function(update) {});
	}

	//   显示点到
	$scope.show_btn = function(item) {
		if(item.show) {
			item.show = false;
		} else {
			for(var i = 0; i < $scope.pointtos.length; i++) {
				$scope.pointtos[i].show = false;
			}
			item.show = true;
		}
	}

	//  点到事件
	$scope.btn_Pointtos = function(item) {
		if(item.State == "0") {
			$scope.postPointtos(item);
		} else {
			$scope.postPointtosClose(item);
		}
	}

	//   点到按钮定位  最右边变朝左边
	var kqjz_position = function() {
		$(document).on("click", ".kqjz-list li", function() {
			//   点到按钮绝对位置left
			var li_left = $(this).position().left;
			var document_width = $(document).width();
			var click_btn = $(this).find(".click-btn");
			if(document_width < (li_left + $(this).width() * 2 + 30)) {
				click_btn.css({
					"left": "-0.95rem"
				});
				click_btn.find("i").css({
					"left": "1.16rem",
					"transform": "rotateY(180deg)"
				});
			}
		});
	}

	//   刷新
	$scope.refresh = function() {
		if($scope.thisCourse.Curriculumclassroomchaptercentreids != undefined && $scope.thisCourse.Curriculumclassroomchaptercentreids != ""){
			$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreids);
		}else{
			$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreid);
		}
	}

	//   排序
	var jz_sort = true;
	$scope.stateSort = "";
	$scope.kqjz_sort = function() {
		if(jz_sort) {
			$scope.stateSort = 'State';
			jz_sort = !jz_sort;
		} else {
			$scope.stateSort = '-State';
			jz_sort = !jz_sort;
		}

	}
	
	//开始定义定时器
	var tm = $scope.setglobaldata.gettimer("zkmb_kqjz");
	if(tm.Key!="zkmb_kqjz"){
		tm.Key="zkmb_kqjz";
		tm.keyctrl="app.zkmb.kqjz";
		tm.fnAutoRefresh=function(){
			this.interval = $interval(function() {
				//   
				if($scope.thisCourse.Curriculumclassroomchaptercentreids != undefined && $scope.thisCourse.Curriculumclassroomchaptercentreids != ""){
					$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreids);
				}else{
					$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreid);
				}
			}, config.zkmb_config.zkmbRefreshTime);	
		};
		tm.fnStopAutoRefresh=function(){
			//console.log("进入取消方法");
			if(!angular.isUndefined(this.interval)) {
				$interval.cancel(this.interval);
				this.interval = 'undefined';
				//console.log("进入取消成功");
			}
			this.interval=null;
		};
		$scope.setglobaldata.addtimer(tm);
	}
	//结束定义定时器
	
	
	//  run
	var run = function(){
		if($scope.thisCourse.Curriculumclassroomchaptercentreids != undefined && $scope.thisCourse.Curriculumclassroomchaptercentreids != ""){
			$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreids);
		}else{
			$scope.getPointtos($scope.thisCourse.Curriculumclassroomchaptercentreid);
		}
		//   定时器
		tm.fnAutoRefreshfn(tm);
	}
	run();

}]);

/*  班级出勤统计     */
app.controller("zkmbKqtjindexContr", ['$scope', '$location', '$rootScope', 'httpService', '$modal', function($scope, $location, $rootScope, httpService, $modal) {

	//
	$scope.xhgl_list = [];

	//  查找叶子
	$scope.isLeaf = function(id) {
		var bol = true;
		for(var i = 0; i < $scope.xhgl_list.length; i++) {
			if($scope.xhgl_list[i].Pld == id) {
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

	//   查看
	$scope.opendet = function(ite) {
		var modalInstance = $modal.open({
			templateUrl: '../project/zkmb/html/zkmb/class/modal_kqtj.html',
			controller: 'qxglKqtjDetailsindexContr',
			windowClass: 'm-modal-kqjz',
			//size: size,
			resolve: {
				items: function() {
					return ite;
				}
			}
		});
	}

	//   取出勤统计
	$scope.getfilterItems = function() {
		var url = config.HttpUrl + "/curriculum/getcurriculumchaptersinfo";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Os": "WEB",
			"Token": config.GetUser().Token,
			"Teacherid": Number(config.GetUser().Usersid),
			"Classesid": 10
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.kqtj_items = data.Result;
				//   课程总出勤率
				for(var a in $scope.kqtj_items){
					var temp = 0;
					for(var b in $scope.kqtj_items[a].Infos){
						temp +=  $scope.kqtj_items[a].Infos[b].Toclassrate;
					}
					$scope.kqtj_items[a].Toclassrate = temp / $scope.kqtj_items[a].Infos.length;
				}
				
			} else {
				$scope.tipService.setMessage(data.Reason);
			}
			console.log('取出勤统计',data)
		}, function(reason) {}, function(update) {});
	}
	
	//   run
	$scope.getfilterItems();

}]);

/*   出勤统计   -  详情   查看   */
app.controller("qxglKqtjDetailsindexContr", ['$scope', 'items', 'httpService','$modalInstance', function($scope, items, httpService,$modalInstance) {
	//console.log("出勤统计---查看详情")

	$scope.item = items;

	$scope.pointtos = [];
	//    显示点到详情
	$scope.getPointtos = function(id) {
		if(!id)return false;
		var myDate = new Date();
		var url = config.HttpUrl + "/action/getpointtos";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
//			"Curriculumclassroomchaptercentreid": Number(id)
			"Curriculumclassroomchaptercentreids":id.toString()//"902,1112"
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.pointtos = data.Result;

				//console.log("取点到详情")
				//console.log(data)
			} else {
				//console.log("取点到详情失败")
				//$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	
	//  close
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};
	
	//   run
	$scope.getPointtos($scope.item.Curriculumclassroomchaptercentreid);
	/*if($scope.item.Curriculumclassroomchaptercentreids != undefined && $scope.item.Curriculumclassroomchaptercentreids != ""){
		$scope.getPointtos($scope.item.Curriculumclassroomchaptercentreids);
	}else{
		$scope.getPointtos($scope.item.Curriculumclassroomchaptercentreid);
	}*/

}]);

/*        设备控制             */
app.controller("zkmbSbkzindexContr", ['$scope', '$state', '$stateParams', '$location', '$rootScope', 'httpService', function($scope, $state, $stateParams, $location, $rootScope, httpService) {

	$scope.sbkz_item;

	var DeviceId = $location.search().DeviceId;
	var DevicePage = $location.search().page;

	//console.log(DevicePage)

	$scope.page = "../project/zkmb/html/zkmb/sbkz/" + DevicePage;
	
	//  GET设备
	var get_info = function(uid, id, type) {
		var url = config.zkmb_config.coapServer + '/device/node/state/device';
		var data = {
			UserID: uid, 		//用户ID(*字符串)   
			DeviceIDs: [id], 	//设备IDs(*字符串数组)    
		    Params: "" 	//参数(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			$scope.sbkz_item = data.Data[0];
		}, function(reason) {}, function(update) {})
	}
	//  run
	get_info(config.zkmb_config.Uid, DeviceId, "device");
}]);