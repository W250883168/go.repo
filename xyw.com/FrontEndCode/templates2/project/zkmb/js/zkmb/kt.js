'use strict';
/**
 * Created by Administrator on 2016/7/29.
 */
/*        设备控制 ---空调            */
app.controller("ktController", ['$scope', '$timeout', 'httpService', '$location','$modal','$filter', function($scope, $timeout, httpService, $location,$modal,$filter) {
	//   空调

	$scope.deviceId = $location.search().DeviceId;
	$scope.classroomId = $location.search().classroomid;
	$scope.deviceimg = config.zkmb_config.deviceImg;
	$scope.uid = config.zkmb_config.Uid;
	
	$scope.model = {
		title: "模式",
		default: {
			code: 0x00,
			name: "自动"
		},
		current: {
			code: 0x00,
			name: "自动"
		},
		list: [{
			code: 0x00,
			name: "自动"
		}, {
			code: 0x01,
			name: "制冷"
		}, {
			code: 0x02,
			name: "除湿"
		}, {
			code: 0x03,
			name: "送风"
		}, {
			code: 0x04,
			name: "制暖"
		}]
	}

	$scope.windSpeed = {
		title: "风速",
		default: {
			code: 0x00,
			name: "自动"
		},
		current: {
			code: 0x00,
			name: "自动"
		},
		list: [{
			code: 0x00,
			name: "自动"
		}, {
			code: 0x01,
			name: "低"
		}, {
			code: 0x02,
			name: "中"
		}, {
			code: 0x03,
			name: "高"
		}]
	}
	$scope.windDirection = {
		title: "风向",
		default: {
			code: 0x00,
			name: "自动"
		},
		current: {
			code: 0x00,
			name: "自动"
		},
		list: [{
			code: 0x00,
			name: "自动"
		}, {
			code: 0x01,
			name: "手动"
		}]
	}
	$scope.temp = {
		title: "温度",
		default: 24,
		current: 24
	}

	//获取设备状态
	$scope.fnGetDeviceStatus = function() {
		//debugger;
		var url = config.zkmb_config.coapServer + "/device/node/state/device";
		var data = {
			UserID: $scope.uid, 		//用户ID(*字符串)   
			DeviceIDs: [$scope.deviceId], 	//设备IDs(*字符串数组)    
		    Params: "" 	//参数(字符串)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			$scope.deviceData = rdata.Data[0];
			//console.log($scope.deviceData);
			$scope.status($scope.deviceData);
			fnInitPara(); //初始化处理
			$.unblockUI();
		}, function(reason) {}, function(update) {});
	}

	//  状态
	$scope.DeviceStatus = {
		//   使用时间
		ServiceTime: 0,
		//   灯泡寿命
		Status: null
	};
	$scope.status = function(item) {
		if(!angular.isObject(item)) return;
		//  使用时间
		$scope.DeviceStatus.ServiceTime = item.UseTimeAfter + item.UseTimeBefore;
		$scope.DeviceStatus.Status = item.DeviceStatus;
	}

	//打开
	$scope.fnOpenDevice = function() {
		//debugger;
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;ִ执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/on/device";
		var data = {
			UserID:$scope.uid,        //用户ID：字符串   
			DeviceID:$scope.deviceId,    //设备ID：字符串   
			Params:""     //参数：字符串(可选)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			fnManualRefresh();
			//$.unblockUI();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

	//关闭
	$scope.fnCloseDevice = function() {
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/off/device";
		var data = {
			UserID:$scope.uid,        //用户ID：字符串   
			DeviceID:$scope.deviceId,    //设备ID：字符串   
			Params:""     //参数：字符串(可选)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			fnManualRefresh();
			//$.unblockUI();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

	//改变模式
	$scope.fnChangeModel = function(index) {
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/controlother";
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "model",
			Para: '"p1":"' + $scope.fnCmdCombind(0x05, $scope.model.list[index].code, 0x08, 0x08) + '"',
			IsSave: "yes", //是否要保存
			SaveValue: index.toString(), //保存的具体值
			IsCreateLog:"yes",                           //modiby by snock:2015-09-23,add parameter,是否创建日志：字符串 空或yes表示创建日志，no表示不创建日志
            AddLogInfo:$scope.model.list[index].name    //modiby by snock:2015-09-23,add parameter,附加日志信息：字符串
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			$scope.model.current = $scope.model.list[index];
			//$.unblockUI();
			fnSendInitBeginCmd(url, data);
			fnManualRefresh();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

	//改变风速
	$scope.fnChangeWindSpeed = function(index) {
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/controlother";
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "windspeed",
			Para: '"p1":"' + $scope.fnCmdCombind(0x07, $scope.windSpeed.list[index].code, 0x08, 0x08) + '"',
			IsSave: "yes", //是否要保存
			SaveValue: index.toString(), //保存的具体值
			IsCreateLog:"yes",                              //modiby by snock:2015-09-23,add parameter, 是否创建日志：字符串 空或yes表示创建日志，no表示不创建日志
            AddLogInfo:$scope.windSpeed.list[index].name    //modiby by snock:2015-09-23,add parameter,附加日志信息：字符串
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			$scope.windSpeed.current = $scope.windSpeed.list[index];
			//$.unblockUI();
			fnSendInitBeginCmd(url, data);
			fnManualRefresh();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

	//改变风向
	$scope.fnChangeWindDirection = function(index) {
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/controlother";
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "winddirection",
			Para: '"p1":"' + $scope.fnCmdCombind(0x08, $scope.windDirection.list[index].code, 0x08, 0x08) + '"',
			IsSave: "yes", //是否要保存
			SaveValue: index.toString(), //保存的具体值
			IsCreateLog:"yes",                                  //modiby by snock:2015-09-23,add parameter,是否创建日志：字符串 空或yes表示创建日志，no表示不创建日志
            AddLogInfo:$scope.windDirection.list[index].name    //modiby by snock:2015-09-23,add parameter,附加日志信息：字符串
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			$scope.windDirection.current = $scope.windDirection.list[index];
			//$.unblockUI();
			fnSendInitBeginCmd(url, data);
			fnManualRefresh();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

	//改变温度
	$scope.fnChangeTemp = function(v) {
		if(($scope.temp.current + v) > 31 || ($scope.temp.current + v) < 16) {
			return
		}
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/controlother";
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "temp",
			Para: '"p1":"' + $scope.fnCmdCombind(0x06, $scope.temp.current + v, 0x08, 0x08) + '"',
			IsSave: "yes", //是否要保存
			SaveValue: ($scope.temp.current + v).toString(), //保存的具体值
			IsCreateLog:"yes",                               //modiby by snock:2015-09-23,add parameter,是否创建日志：字符串 空或yes表示创建日志，no表示不创建日志
            AddLogInfo:($scope.temp.current+v).toString()    //modiby by snock:2015-09-23,add parameter,附加日志信息：字符串
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			$scope.temp.current = $scope.temp.current + v;
			//$.unblockUI();
			fnSendInitBeginCmd(url, data);
			fnManualRefresh();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

	//空调核码
	$scope.fnCheckCode = function() {
		if($scope.checkCode.code == "") return;
		fnStopAutoRefresh();
		//将空调码转换为十六进制(2字节），并取出高位和低位
		var code = $scope.fnCalcACCode($scope.checkCode.code);
		b2 = parseInt(code.substr(0, 2), 16);
		b3 = parseInt(code.substr(2, 2), 16);

		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/controlother";
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "checkcode",
			Para: '"p1":"' + $scope.fnCmdCombind(0x02, b2, b3, 0x08, 0x08) + '"'
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			val = parseInt($scope.checkCode.code) + 1;
			str = '000' + val.toString();
			$scope.checkCode.code = str.substring(str.length - 3, str.length);

			//发送开机命令（如果听到响声，则表示有效）
			$scope.fnOpenDevice();
			
			fnManualRefresh();
		}, function(reason) {}, function(update) {})
	}
	
	var timer = "";
	//返回
	$scope.fnBack = function() {
		$window.history.back();
	}
	
	//停止自动刷新
	var fnStopAutoRefresh = function() {
		if(!angular.isUndefined(timer)) {
			$timeout.cancel(timer);
			timer = 'undefined';
		}
	}

	//手动刷新
	var fnManualRefresh = function() {
		timer = $timeout(function() {
			$scope.fnGetDeviceStatus();
		}, config.zkmb_config.immediatelyRefreshTime);
	}

	//自动刷新
	var fnAutoRefresh = function() {
		timer = $timeout(function() {
			$scope.fnGetDeviceStatus();
		}, config.zkmb_config.fixRefreshTime);
	}
	
	

	$scope.fnInit = function() {
			$scope.fnGetDeviceStatus();
		}
		//  run
	$scope.fnInit();

	//命令编联
	$scope.fnCmdCombind = function(b1, b2, b3, b4) {
		return "hex" + $scope.b2a(b1) + $scope.b2a(b2) + $scope.b2a(b3) + $scope.b2a(b4) + $scope.b2a(b1 ^ b2 ^ b3 ^ b4);
	}

	//十六进制整数到十六进制字符串(oxFF > FF)
	$scope.b2a = function(b) {
		var tmp = b.toString(16);
		if(tmp.length == 1) {
			tmp = "0" + tmp;
		} else {
			tmp = "" + tmp;
		}
		//将第5个字节拼接到hex后面
		return tmp;
	}

	$scope.fnCalcACCode = function(current) {
		var tmp = parseInt(current).toString(16);
		str = '0000' + tmp;
		return str.substring(str.length - 4, str.length);
	}

	//初始化界面
	function fnInitPara() {
		var sendContents = $scope.deviceData.LastSendContent;
		for(var i = 0; i < sendContents.length; i++) {
			var sendContent = sendContents[i];
			//初始化空调模式
			if(sendContent.CmdCode == "model" && sendContent.Value != "") {
				$scope.model.current = $scope.model.list[parseInt(sendContent.Value)];
			}
			//初始化风速
			if(sendContent.CmdCode == "windspeed" && sendContent.Value != "") {
				$scope.windSpeed.current = $scope.windSpeed.list[parseInt(sendContent.Value)];
			}
			//初始化风向
			if(sendContent.CmdCode == "winddirection" && sendContent.Value != "") {
				$scope.windDirection.current = $scope.windDirection.list[parseInt(sendContent.Value)];
			}
			//初始化温度
			if(sendContent.CmdCode == "temp" && sendContent.Value != "") {
				$scope.temp.current = parseInt(sendContent.Value);
			}
		}
	}

	//----------------------------------------------------------------------------------------------------------------------------
	//-当模式、风速、风向、温度等改变时，将当前参数写入到节点模块内部，下次发送开机命令时候，空调就可以工作在这组参数上了-----------------------
	//发送初始化启用命令
	function fnSendInitBeginCmd(url, cmddata) {
		//modiby by snock:2015-09-23,add the follow code,将IsCreateLog赋值为"no",不创建日志---------------
        cmddata.IsCreateLog = "no";
        
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "initbegin",
			Para: "", //具体的命令配置在后面DB中，所以不需要传入
			IsCreateLog:"no"                               //modiby by snock:2015-09-23,add parameter,是否创建日志：字符串 空或yes表示创建日志，no表示不创建日志
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {
			//发送初始化启用命令后，发送初始化命令
			fnSendInitCmd(url, cmddata)
		}, function(reason) {}, function(update) {})
	}

	//发送初始化命令
	function fnSendInitCmd(url, cmddata) {
		var promise = httpService.ajaxPost(url, cmddata);
		promise.then(function(rdata) {
			//发送初始化命令后，发送初始化结束命令
			fnSendInitEndCmd(url)
		}, function(reason) {}, function(update) {})
	}

	//发送初始化结束命令
	function fnSendInitEndCmd(url) {
		var data = {
			UId: $scope.uid,
			DeviceId: $scope.deviceId,
			UseWhoseCmd: "self",
			CmdCode: "initend",
			Para: "", //具体的命令配置在后面DB中
			IsCreateLog:"no"                               //modiby by snock:2015-09-23,add parameter,是否创建日志：字符串 空或yes表示创建日志，no表示不创建日志
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(rdata) {}, function(reason) {}, function(update) {})
	}

	//   取消定时器
	$scope.$on("$destroy", function() {
		//$scope.zkmb_fnGetDeviceStatus();
	});
	
	
	
	
	
	/*--------------------------  设备数据   ------------------------------*/
	//   默认第一页
	var PageIndex = 1;
	//   分页页显示条数
	var PageSize = 9;
	//   显示页码数
	var pageNumber = 5;
	//  设备预警
	$scope.AlertInfoList = [];
	//  设备故障
	$scope.FaultInfoList = [];
	//  设备操作日志
	$scope.OperateLogList = [];
	//  设备使用日志
	$scope.UseLogList = [];
	//  分页
	$scope.backPage2 = {};
	$scope.backPage3 = {};
	$scope.backPage4 = {};
	$scope.backPage5 = {};

	//  设备预警
	var getAlertInfoList = function(pageindex) {
		var url = config.HttpUrl + "/device/getAlertInfoList";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				DeviceId: $scope.deviceId
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.AlertInfoList = data.Result.Data;
				//   分页
				//$scope.backPage2 = pageFn(data.Result.Page,pageNumber);
				console.log("预警")
				//console.log($scope.AlertInfoList)
			} else {
				alert(data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//  设备故障
	var getFaultInfoList = function(pageindex) {
		var url = config.HttpUrl + "/device/getFaultInfoList";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: pageindex, 
				PageSize: 8
			},
			Para: {
				DeviceId: $scope.deviceId
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.FaultInfoList = data.Result.Data;
				//   分页
				$scope.backPage3 = pageFn(data.Result.Page,pageNumber);
			} else {
				alert(data.Reason);
			}
			console.log("设备故障")
			console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	//  设备操作日志
	var getOperateLogList = function(pageindex) {
		var url = config.HttpUrl + "/device/getOperateLogList";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: pageindex, 
				PageSize: PageSize
			},
			Para: {
				DeviceId: $scope.deviceId
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.OperateLogList = data.Result.Data;
				//   分页
				$scope.backPage4 = pageFn(data.Result.Page,pageNumber);
			} else {
				alert(data.Reason);
			}
			console.log("设备操作日志")
			console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	//  设备使用日志
	var getUseLogList = function(pageindex) {
		var url = config.HttpUrl + "/device/getUseLogList";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: pageindex, 
				PageSize: PageSize
			},
			Para: {
				DeviceId: $scope.deviceId
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.UseLogList = data.Result.Data;
				//   分页
				$scope.backPage5 = pageFn(data.Result.Page,pageNumber);
				//   加入时长
				for(var i = 0; i < $scope.UseLogList.length; i++){
					var offtime = Date.parse(new Date($scope.UseLogList[i].OffTime.replace(/\-/g, "/")));
					var ontime = Date.parse(new Date($scope.UseLogList[i].OnTime.replace(/\-/g, "/")));
					if(!offtime && !(offtime > ontime)){
						$scope.UseLogList[i].timeLength = 0;
					}else{
						var lengthtime = $filter('FormatTime')((offtime - ontime) / 1000);
						$scope.UseLogList[i].timeLength = lengthtime;
					}
					
				}
			} else {
				alert(data.Reason);
			}
			console.log("设备使用日志")
			console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	
	
	/*  -------------------- 分页、页码  -----------------------  */
	/*----------------
	//    分页对象添加页码
	//    return  obj  分页对象
	//    pagedata:obj  分页对象
	//    maxpagenumber:int  显示页码数默认5个页码
	------------------*/
	var pageFn = function(pagedata,maxpagenumber){
		if(pagedata.length < 1)return null;
		//   缺省时分5页
		Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
		var nub = [];
		var mid = Math.ceil(maxpagenumber / 2);
		if(pagedata.PageCount > maxpagenumber){
			//  起始页
			var Snumber = 1;
			if((pagedata.PageIndex - mid) < 1 ){
				Snumber = 1
			}else if((pagedata.PageIndex + mid) > pagedata.PageCount){
				Snumber = pagedata.PageCount - maxpagenumber + 1;
			}else{
				Snumber = pagedata.PageIndex - (mid - 1)
			}
			for(var i = 0; i < maxpagenumber; i++){
				nub.push(Snumber + i);
			}
		}else{
			for(var i = 0; i < pagedata.PageCount; i++){
				nub.push(i + 1);
			}
		}
		pagedata.Number = nub;
		return pagedata;
	}
	//  翻页
	$scope.pageClick = function(tabid,pageindex){
		Number(pageindex) < 1 ? pageindex = 1 : pageindex = Number(pageindex);
		if(pageindex > $scope['backPage' + tabid].PageCount)return false;
		switch(tabid){
			//   设备故障
			case 3:getFaultInfoList(pageindex); break;
			//   设备操作日志
			case 4:getOperateLogList(pageindex); break;
			//   使用日志
			case 5:getUseLogList(pageindex); break;
		}
	}
	/*  -------------------- 分页、页码  -----------------------  */
	
	/*   ------------------  添加，查看，维护   --------------------   */   
	
	
	/* -------- 添加  -------- */
	
	//
	$scope.form = {
		"add":{}
	}
	//   
	$scope.form.add = {
		//   故障id：字符
		"Id":"",
		//   故障设备
		"DeviceId":"",
		//   故障设备名称
		"DeviceName":"",
		//   教室所有设备Items
		"DeviceItems":[],
		//   设备位置
		"DeviceSite":"",
		//   故障现象
		"FaultSummary":"",
		//   故障描述
		"FaultDescription":"",
		//   故障发生时间
		"HappenTime":"",
		//   设备是否可用  0/1(0-不可使用 1-可以使用)
		"IsCanUse":"0",
		//   申报人 id
		"InputUserId":"",
		//   申报人 名称
		"InputUserName":"",
		//   申报时间
		"InputTime":"",
		//   提交时间（提交故障时间）
		"SubmitTime":"",
		//   故障状态 字符（0-草稿 1-待受理 2-维修中 3-已维修）
		"Status":"0"
	}
	
	//   删除故障
	var deleteFault = function(id){
		if(!id)return false;
		var url = config.HttpUrl + "/device/deleteFault";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				//故障id：字符，不能为空
				Id:id         
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				alert("删除成功！");
				getFaultInfoList($scope.backPage3.PageIndex);
			} else {
				alert(data.Reason);
			}
			console.log("删除故障")
			console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	//   故障添加打开弹窗
	$scope.modalOpenAddFault = function(item,str) {
		if(str == 'delete'){
			if(!confirm("确定删除？")) return false;
			deleteFault(item.Id);
			return false;
		}
		
		if(!item){
			item = {
				"deviceId":$scope.deviceId,
				"classroomId":$scope.classroomId,
				"DeviceSite":$scope.classRoomItem.Campusname + "-" + $scope.classRoomItem.Buildingname + "-" + $scope.classRoomItem.Floorname + "-" + $scope.classRoomItem.Classroomsname
			}
		}else{
			item.classroomId = $scope.classroomId;
		}
		
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_device_fault.html',
			controller: 'modalFaultAddCtrl',
			resolve: {
				items: function() {
					return [item,str];
				},
				deps: ['$ocLazyLoad',
					function($ocLazyLoad) {
						return $ocLazyLoad.load(['ui.select']);
					}
				]
			}
		});
		
		modalInstance.result.then(function(bul) {
			console.log(bul)
			if(bul){
				getFaultInfoList($scope.backPage3.PageIndex);
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	
	
	
	
	/* -------- 添加 End -------- */
	
	
	
	
	/*   ------------------  添加，查看，维护 End   --------------------   */   
	

	//   run
	var info_run = function(){
		//getClassroomInfo($scope.classroomId);
		getAlertInfoList();
		getFaultInfoList(PageIndex);
		getOperateLogList(PageIndex);
		getUseLogList(PageIndex);
	}
	info_run();
	//  
	
	//   分页
	
	/*--------------------------  设备数据 End   ------------------------------*/

}]);