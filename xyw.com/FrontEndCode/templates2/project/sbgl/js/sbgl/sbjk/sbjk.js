'use strict';
/**
 * Created by Administrator on 2016/7/21.
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
/*   设备监控    */
app.controller('sbjkContr', ['$scope', 'httpService', '$modal', '$interval', '$rootScope','$location','toaster', function($scope, httpService, $modal, $interval, $rootScope,$location,toaster) {
	//$scope.alerts.push({type: 'success',msg: '开启关怀成功!'});
	console.log("设备监控")

	//   查看楼层平面图
	$scope.openPic = function(floor) {
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbjk/modal_pic.html',
			controller: 'picSbjkContr',
			windowClass: 'm-sbjk-modal',
			//size: "lg",
			resolve: {
				items: function() {
					return floor;
				}
			}
		});
	}

	//   停用 预警 离线
	$scope.anomalyExtension = {
		stop: false,
		warning: true,
		offline: true
	}


	//   默认校区  教学楼  本地存储
	$scope.defaultSchoolFloor = {
		"school": [],
		"floor": []
	};
	//开始定义定时器
	var tm=$scope.setglobaldata.gettimer("sbjk");
	if(tm.Key!="sbjk"){
		tm.Key="sbjk";
		tm.keyctrl="app.sbgl.sbjk";
		tm.fnAutoRefresh=function(){
			console.log("开始调用定时器");
			this.interval = $interval(function() {
				buildingString();//   格式化字符串
				getClassroomStatusList(buildingids);//   get 教室
			}, config.sbjkRefreshTime);
		};
		tm.fnStopAutoRefresh=function(){
			console.log("进入取消方法");
			if(!angular.isUndefined(this.interval)) {
				$interval.cancel(this.interval);
				this.interval = 'undefined';
				console.log("进入取消成功");
			}
			this.interval=null;
		};
		$scope.setglobaldata.addtimer(tm);
	}
	//结束定义定时器

	//   取校区
	$scope.schoolItems = [];
	var getcampus = function() {
		var url = config.HttpUrl + "/basicset/getcampus";
		var data = {

		};
		var promise = httpService.ajaxGet(url, null);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				//   自动刷新
				//var tm=$scope.setglobaldata.gettimer("sbjk");
				tm.fnAutoRefreshfn(tm);
				//AutoRefresh();

				$scope.schoolItems = data.Result;
				//   加入默认选中项
				if($scope.defaultSchoolFloor["school"].length > 0) {
					for(var i = 0; i < $scope.schoolItems.length; i++) {
						if($scope.schoolItems[i].Campusid == $scope.defaultSchoolFloor["school"][0].Campusid) {
							$scope.schoolItems[i].checkbox = true;
						}
					}
					getcampusFloor($scope.defaultSchoolFloor["school"][0].Campusid);
				} else {
					$scope.schoolItems[0].checkbox = true;
					getcampusFloor($scope.schoolItems[0].Campusid);
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}
//	getcampus();

	//  取楼栋
	$scope.floorItems = [];
	var getcampusFloor = function(campusid) {
		var url = config.HttpUrl + "/basicset/getbuilding";
		var data = {
			//"Usersid": config.GetUser().Usersid,
			//"Rolestype": config.GetUser().Rolestype,
			//"Token": config.GetUser().Token,
			//"Os": "WEB",
			"campusid": campusid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.floorItems = data.Result;
				//   加入默认选中项
				for(var i = 0; i < $scope.floorItems.length; i++) {
					if($scope.defaultSchoolFloor["floor"].length>0){
					for(var b = 0; b < $scope.defaultSchoolFloor["floor"].length; b++) {
						if($scope.floorItems[i].Buildingid == $scope.defaultSchoolFloor["floor"][b].Buildingid) {
							$scope.floorItems[i].checkbox = true;
						}
					}
					}else{
						$scope.floorItems[i].checkbox = true;
					}
				}
				if($scope.defaultSchoolFloor["floor"].length==0){
				$scope.floor_checkbox();
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}

	//   校区选择
	$scope.school_tab = function(item) {
		for(var i = 0; i < $scope.schoolItems.length; i++) {
			if($scope.schoolItems[i].Campusid == item.Campusid) {
				$scope.schoolItems[i].checkbox = true;
			} else {
				$scope.schoolItems[i].checkbox = false;
			}
		}
		//清除楼栋选择数据
		$scope.defaultSchoolFloor["floor"]=[];
		if(localStorage.getItem("sbjkTab"+item.Campusid)!=undefined){
		var defaultSchoolFloorlib = JSON.parse(localStorage.getItem("sbjkTab"+item.Campusid));
		$scope.defaultSchoolFloor=defaultSchoolFloorlib;
		}
		//   查楼栋
		getcampusFloor(item.Campusid);

		if($scope.defaultSchoolFloor.school[0] != undefined){
			if(item.Campusid == $scope.defaultSchoolFloor.school[0].Campusid){
				//   格式化字符串
				buildingString();
				//   get 教室
				getClassroomStatusList(buildingids);
			}else{
				$scope.classroomItems = [];
			}
		}
	}

	//   楼栋选择
	$scope.floor_checkbox = function(item) {

		//   校区
		$scope.defaultSchoolFloor["school"] = [];
		for(var i = 0; i < $scope.schoolItems.length; i++) {
			if($scope.schoolItems[i].checkbox) {
				$scope.defaultSchoolFloor["school"].push($scope.schoolItems[i]);
			}
		}
		//   楼栋
		$scope.defaultSchoolFloor["floor"] = [];
		for(var i = 0; i < $scope.floorItems.length; i++) {
			if($scope.floorItems[i].checkbox) {
				$scope.defaultSchoolFloor["floor"].push($scope.floorItems[i]);
			}
		}

		//  格式化
		var temp = null;
		temp = $scope.defaultSchoolFloor["school"][0];
		$scope.defaultSchoolFloor["school"][0] = {};
		$scope.defaultSchoolFloor["school"][0].Campusid = temp.Campusid;
		$scope.defaultSchoolFloor["school"][0].Campusname = temp.Campusname;
		for(var i = 0; i < $scope.defaultSchoolFloor["floor"].length; i++) {
			temp = $scope.defaultSchoolFloor["floor"][i];
			$scope.defaultSchoolFloor["floor"][i] = {};
			$scope.defaultSchoolFloor["floor"][i].Buildingid = temp.Buildingid;
			$scope.defaultSchoolFloor["floor"][i].Buildingname = temp.Buildingname;
		}
		localStorage.setItem("sbjkTab"+$scope.defaultSchoolFloor["school"][0].Campusid, JSON.stringify($scope.defaultSchoolFloor));
		localStorage.setItem("sbjkTab", JSON.stringify($scope.defaultSchoolFloor));
//console.log(localStorage.getItem("sbjkTab"));

		//   格式化字符串
		buildingString();
		//   get 教室
		getClassroomStatusList(buildingids);
	}

	/////////////////////////////////////////////////////////

	//------------------------------------------------------------------------------------------------------------------
	var fnTransformData = function(data) {
		var d;
		d = data;

		//定义变量
		var bid = undefined; //楼栋id
		var fid = undefined; //楼层id
		var ob = undefined; //object of building
		var of = undefined; //object of floor
		var oc = undefined; //object of classroom
		var xh = 0; //序号
		var maxCol = 0; //所有楼层中的最大教室数量（所有楼层按这个数量显示教室列数）
		var rd = []; //存放最后的数据

		//循环对数据进行处理,构造更易在界面上展现的数据集
		for(var i = 0; i < d.length; i++) {
			var ir2 = null;
			xh++; //序号直接加1(遇新楼层时，恢复为1

			//取出教室编号后两位
			var name = d[i].ClassroomName;
			ir2 = parseInt(name.substring(name.length - 2, name.length)); // 将编号右边两位转换为整数

			//oc
			oc = {
				ClassroomId: d[i].ClassroomId,
				ClassroomName: d[i].ClassroomName,
				ClassroomState: d[i].ClassroomState,
				CollectionNumbers: d[i].CollectionNumbers,
				HaveStop: d[i].HaveStop,
				HaveAlert: d[i].HaveAlert,
				HaveOffline: d[i].HaveOffline,
				HaveRun: d[i].HaveRun
			}

			//of
			if(d[i].FloorId != fid) {
				//将当前楼层id保存到fid
				fid = d[i].FloorId;

				//将前一个楼层的最大教室数保存起来
				if(xh - 1 > maxCol) {
					maxCol = xh - 1
				}

				//保存前一个
				if(of != undefined) {
					ob.data.push(of);
					of = undefined;
				}

				//建立新的楼层
				of = {
					FloorId: d[i].FloorId,
					FloorName: d[i].FloorName,
					FloorImage: d[i].FloorImage,
					data: []
				};
				xh = 1; //建立新的楼层后，教室序号肯定是从1开始，所以初始化为1
			}

			//将教室压入楼层前,先补插空缺教室（按顺序补齐后台数据中没有返回的教室）
			for(var k = 0; k < ir2 - xh; k++) {
				of.data.push(fnGetNullClassroom(name.substring(0, name.length - 2), xh)) //
				xh++; //将序号+1
			}

			//将教室压入楼层
			of.data.push(oc)

			//ob
			if(d[i].BuildingId != bid) {
				//保存前一个
				if(ob != undefined) {
					rd.push(ob); //将ob压入rd
					ob = undefined;
				}
				bid = d[i].BuildingId;

				//建立新的楼栋
				ob = {
					BuildingId: d[i].BuildingId,
					BuildingName: d[i].BuildingName,
					data: []
				};
			}

		}

		//最后一栋
		if(ob != undefined) {
			if(of != undefined) {
				ob.data.push(of) //最后一层压入楼栋
			}
			rd.push(ob); //将ob压入rd
		}

		//将前一个楼层的最大教室数保存起来
		if(xh > maxCol) {
			maxCol = xh
		}

		//返回
		return {
			MaxCol: maxCol,
			data: rd
		}
	}

	//获得空教室对象
	var fnGetNullClassroom = function(floorCode, classroomName) {
		//将教室编码补齐
		var fullClassroomName
		if(classroomName < 10) {
			fullClassroomName = floorCode + "0" + classroomName
		} else {
			fullClassroomName = floorCode + classroomName
		}
		//凡是补插的空教室，教室id一律为-1
		return {
			ClassroomId: -1,
			ClassroomName: fullClassroomName,
			ClassroomState: -1,
			CollectionNumbers: 0,
			HaveStop: -1,
			HaveAlert: -1,
			HaveOffline: -1,
			HaveRun: -1
		}
	}

	//增加楼层教室（补齐楼层教室）
	var fnAddFloorClassroom = function(data) {
			//以所有楼层最大教室数量为参考，将楼层教室数量小于最大数量的，补齐
			var rd = data.data;
			var maxCol = data.MaxCol;
			for(var i = 0; i < rd.length; i++) {
				var f = rd[i].data;
				for(var j = 0; j < f.length; j++) {
					var c = f[j].data;
					//如果当前楼层的教室数量小于maxCol，则循环补齐
					if(c.length < maxCol) {
						//得到教室名称的前缀(使用第一个教室的名称，除去后两位即可)
						name = c[0].ClassroomName;
						//当前楼层已有教室数量
						var len = c.length;
						//循环补齐
						for(var k = 0; k < maxCol - len; k++) {
							c.push(fnGetNullClassroom(name.substring(0, name.length - 2), len + 1 + k))
						}
					}
				}
			}

			return {
				MaxCol: maxCol,
				data: rd
			};
		}
		//------------------------------------------------------------------------------------------------------------------
		/////////////////////
	//   取教室列表
	//  被选中楼栋ID
	var buildingids = "";
	var buildingString = function(){
		//   楼栋
		buildingids = "";
		for(var i = 0; i < $scope.defaultSchoolFloor.floor.length; i++) {
			//   楼栋ID
			buildingids += $scope.defaultSchoolFloor.floor[i].Buildingid + ",";
		}
		if(buildingids.length > 0) {
			buildingids = buildingids.substr(0, buildingids.length - 1);
		}else{
			buildingids = "";
		}
	}



	//   设备控制 取教室列表
	$scope.classroomItems = [];
	var getClassroomStatusList = function(buildingids) {
		if(!buildingids){$scope.classroomItems = [];return false;}
		var url = config.HttpUrl + "/device/getClassroomStatusList";
		var data = {
			Auth: {
				//用户id：整型
				"Usersid": config.GetUser().Usersid,
				//角色类型：整型
				"Rolestype": config.GetUser().Rolestype,
				//令牌：字符串
				"Token": config.GetUser().Token,
				//操作系统：字符串
				"Os": "WEB"
			},
			Para: {
				//多个楼栋id，逗号拼接成串
				BuildingIds: buildingids
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.classroomItems = fnTransformData(data.Result.Data).data;
				//console.log(1,$scope.classroomItems)
				//$scope.classroomItems = fnAddFloorClassroom($scope.classroomItems).data;
				//console.log(2,$scope.classroomItems)

			} else {
				$scope.classroomItems = [];
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//   楼层一键关闭
	$scope.tierOff = function(item) {
		console.log(item)
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/sbgl/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/off/floor";
		var data = {
			"UserID": config.GetUser().Usersid.toString(),	//(*字符串)
			"FloorID": item.FloorId.toString(), //(*字符串)
			"Params": ""	//(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {} else {
				$.unblockUI();
			}
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}


	//   载入默认教室
	if(localStorage.getItem("sbjkTab") != null) {
		$scope.defaultSchoolFloor = JSON.parse(localStorage.getItem("sbjkTab"));
		//   格式化字符串
		buildingString();
		//   get 教室
		getClassroomStatusList(buildingids);
	}else{

		}

	$scope.$on('to-parent', function(event,data) {
        console.log('ParentCtrl', data);       //父级能得到值
   });


	//  run
	var run = function() {
		getcampus();
	}
	run();


}]);

/*   设备监控列表页 教室设备列表     */
app.controller('sbjkListContr', ['$scope', 'httpService', '$stateParams', '$interval','$state','$filter','$rootScope','toaster', function($scope, httpService, $stateParams, $interval,$state,$filter,$rootScope,toaster) {
	console.log($stateParams.classroomid)
	//  timer
	var timer = null;
	//  教室ID
	$scope.classroomid = $stateParams.classroomid;
	//$scope.classroomid = "104";
	//  设备数组
	$scope.deviceData = [];
	//  图片地址
	$scope.deviceimg = config.zkmb_config.deviceImg;
	//   教室信息
	$scope.classRoomItem = {};
	//   点击某个设备时的设备名称  用于面包屑
	$scope.deviceItemName = "";


	//开始定义定时器
	var tm=$scope.setglobaldata.gettimer("sbjk_list");
	if(tm.Key!="sbjk_list"){
		tm.Key="sbjk_list";
		tm.keyctrl="app.sbgl.sbjk.list";
		tm.fnAutoRefresh=function(){
			//console.log("开始调用定时器");
			this.interval = $interval(function() {
				$scope.zkmb_fnGetDeviceStatus();
				//   取教室信息
				$scope.getClassroomInfo($scope.classroomid);
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

	//    取教室信息
	$scope.getClassroomInfo = function(classroomid) {
		if(Number(classroomid) < 0 && !classroomid){return false}else{classroomid = Number(classroomid)};
		var url = config.HttpUrl + "/basicset/getclassroominfo";
		var data = {
			id: classroomid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.classRoomItem = data.Result;
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	$scope.getClassroomInfo($scope.classroomid);



	//获取设备状态信息
	$scope.zkmb_fnGetDeviceStatus = function() {
		var url = config.zkmb_config.coapServer + "/device/node/state/room";
		var data = {
			UserID: config.GetUser().Usersid.toString(), 		//用户ID(*字符串)
			RoomID: $scope.classroomid, 	//房间ID(*字符串)
		    Params: "" 	//参数(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			$scope.deviceData = data.Data;
			//console.log($scope.deviceData)
			//   自动刷新
			tm.fnAutoRefreshfn(tm);
			$.unblockUI();
		}, function(reason) {}, function(update) {})
	}
		//  run
	$scope.zkmb_fnGetDeviceStatus();

	//   教室设备 一键开启
	$scope.fnOpenDevice = function() {
		tm.fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/sbgl/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/on/room";
		var data = {
			"UserID": config.GetUser().Usersid.toString(),	//(*字符串)
			"RoomID": $scope.classroomid, //(*字符串)
			"Params": ""	//(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			tm.fnAutoRefresh();
		}, function(reason) {}, function(update) {})
	}

	//教室设备 一键关闭
	$scope.fnCloseDevice = function() {
		tm.fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/sbgl/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/off/room";
		var data = {
			"UserID": config.GetUser().Usersid.toString(),	//(*字符串)
			"RoomID": $scope.classroomid, //(*字符串)
			"Params": ""	//(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			tm.fnAutoRefresh();
		}, function(reason) {}, function(update) {})
	}

	/*//停止自动刷新
	var fnStopAutoRefresh = function() {
		if(!angular.isUndefined(timer)) {
			$timeout.cancel(timer);
			timer = 'undefined';
		}
	}*/

	/*//手动刷新
	var fnManualRefresh = function() {
		$timeout(function() {
			$scope.zkmb_fnGetDeviceStatus();
		}, config.zkmb_config.immediatelyRefreshTime);
	}

	//自动刷新
	var fnAutoRefresh = function() {
		timer = $timeout(function() {
			$scope.zkmb_fnGetDeviceStatus();
		}, config.zkmb_config.fixRefreshTime);
	}

	//   后退离开停止定时器
	$scope.$on("$destroy", function() {
		fnStopAutoRefresh();
	});*/

	//   点击打开设备
	$scope.openDevice = function(item){
		//    设备名称用于面包屑
		for(var i = 0; i < $scope.deviceData.length; i++){
			if($scope.deviceData[i].DeviceId == item.DeviceId){
				$scope.deviceItemName = item.DeviceName;
				break;
			}
		}
		$state.go("app.sbgl.sbjk.list.article",{'DeviceId':item.DeviceId,'page':item.DevicePage});
	}

	/*//   监听离开页面取消定时器
	$rootScope.$on('$stateChangeStart',
		function(event, toState, toParams, fromState, fromParams) {
			if(toState.name != "app.sbgl.sbjk.list") {
				fnStopAutoRefresh();
			} else if(fromState.name == "app.sbgl.sbjk.list.article") {
				//  开启自动刷新
				fnAutoRefresh();
			}
		}
	);
*/
}]);

/*   设备监控列表页 教室设备详情      */
app.controller('sbjkArticleContr', ['$scope', '$location', function($scope, $location) {
	//console.log($location)
	var ClassroomId = $location.search().classroomid;
	var DeviceId = $location.search().DeviceId;
	var DevicePage = $location.search().page;
	console.log("教室详情")

	$scope.page = "../project/zkmb/html/zkmb/sbkz/" + DevicePage;

	//   点击某个设备时的设备名称  用于面包屑
	$scope.deviceItemName = "";

	//   监听请求返回
	$scope.$watch('deviceData',function(newValue,oldValue, scope){
		//    设备名称用于面包屑
		for(var a in newValue){
			if(newValue[a].DeviceId == DeviceId){
				$scope.deviceItemName = newValue[a].DeviceName;
				break;
			}
		}
	});

}]);



/*   设备监控列表-查课表       */
app.controller('sbjkListCkbContr', ['$scope', 'httpService', '$location', '$filter','toaster', function($scope, httpService, $location, $filter,toaster) {
	$scope.ClassroomId = $location.search().classroomid;

	//   时间表
	$scope.dateTable = config.dateTable;
	//   内容
	$scope.dateTableBody = {};
	//   星期日期数组
	$scope.dataWeek = {};
	/*  -----------------------------------------------------  */
	/**
	 * 获取本周、本季度、本月、上月的开端日期、停止日期
	 */
	function getMonDate() {
		var d = new Date(),
			day = d.getDay(),
			date = d.getDate();
		if(day == 1)
			return d;
		if(day == 0)
			d.setDate(date - 6);
		else
			d.setDate(date - day + 1);
		return d;
	}
	// 0-6转换成中文名称
	function getDayName(day) {
		var day = parseInt(day);
		if(isNaN(day) || day < 0 || day > 6)
			return false;
		var weekday = ["星期天", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"];
		return weekday[day];
	}
	// d是当前星期一的日期对象
	var d = getMonDate();
	var arr = [];
	for(var i = 0; i < 7; i++) {
		var dd,mo;
		d.getDate() < 10 ? dd = "0" + d.getDate() : dd = d.getDate();
		(d.getMonth() + 1) < 10 ? mo = "0" + (d.getMonth() + 1) : mo = (d.getMonth() + 1);
		arr.push({'name':getDayName(d.getDay()),"date":d.getFullYear() + '-' + mo + '-' + dd});
		d.setDate(d.getDate() + 1);
	}
	$scope.dataWeek = arr;

	/*  -----------------------------------------------------  */


	//   教室信息
	$scope.classRoomItem = {};

	//    取教室信息
	$scope.getClassroomInfo = function(classroomid) {
		if(Number(classroomid) < 0 && !classroomid){return false}else{classroomid = Number(classroomid)};
		var url = config.HttpUrl + "/basicset/getclassroominfo";
		var data = {
			id: classroomid
		};
		var promise = httpService.ajaxGet(url, data);
		promise.then(function(data) {
			//console.log(data)
			if(data.Rcode == "1000") {
				$scope.classRoomItem = data.Result;
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	$scope.getClassroomInfo($scope.classroomid);


	//    查课表时间 //   格式化时间课表
	var tableDate = function(classitems){
		//    课节
		for(var c in $scope.dateTable){
			//   行
			$scope.dateTableBody[c] = {};
			//   日期星期几
			for(var i = 0; i <  $scope.dataWeek.length; i++){
				$scope.dateTableBody[c][$scope.dataWeek[i].date] = {};
				for(var b = 0; b < classitems.length; b++){
					if(classitems[b].Begindate.indexOf($scope.dataWeek[i].date) > -1 && classitems[b].Begindate.indexOf($scope.dateTable[c].substr(0,4)) > -1){
						$scope.dateTableBody[c][$scope.dataWeek[i].date] = classitems[b];
						break;
					}
				}
			}
		}
	}


	//   课表
	$scope.schedule = [];
	//  用户查课表  当天内
	var getCurriculums = function() {
		var myDate = new Date();
		var url = config.HttpUrl + "/action/getcurriculums";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Begindate": arr[0].date + " 00:00:01",
			"Enddate": arr[6].date + " 23:59:59",
			"State": -1,
			//"Teacherids":config.GetUser().Usersid.toString(),
			"Teacherids":"",
			"Classroomid":Number($scope.ClassroomId),
			//"PageSize":70,
			"PageIndex":-1
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("教室查课表 -当天",data)
			if(data.Rcode == "1000") {
				$scope.schedule = data.Result;


				tableDate($scope.schedule);

			} else {
				$scope.tipService.setMessage(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	getCurriculums();

	//   监听请求返回
	$scope.$watch('schedule',function(newValue,oldValue, scope){
		//    设置正在上课 lessoning = true
		for(var a in $scope.schedule){
			if($scope.schedule[a].Curriculumclassroomchaptercentreid == $scope.classRoomItem.Curriculumclassroomchaptercentreid){
				$scope.schedule[a].lessoning = true;
			}
		}
	});
	//   监听请求返回
	$scope.$watch('classRoomItem',function(newValue,oldValue, scope){
		//    设置正在上课 lessoning = true
		for(var a in $scope.schedule){
			if($scope.schedule[a].Curriculumclassroomchaptercentreid == $scope.classRoomItem.Curriculumclassroomchaptercentreid){
				$scope.schedule[a].lessoning = true;
			}
		}
	});
}]);
