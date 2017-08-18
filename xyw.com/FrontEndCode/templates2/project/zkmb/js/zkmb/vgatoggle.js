/**
 * Created by Administrator on 2016/7/29.
 */
app.controller('vgatoggleController',['$scope', 'httpService','$location','$window','$interval','$timeout','$modal','$filter', function ($scope, httpService,$location,$window,$interval,$timeout,$modal,$filter) {
	
	$scope.deviceId = $location.search().DeviceId;
	$scope.classroomId = $location.search().classroomid;
	$scope.deviceimg = config.zkmb_config.deviceImg;
	$scope.uid = config.zkmb_config.Uid;
	
    $scope.port = {
        title:"切换端口",
        default: {},
        //  选中项
        current: {},
        list: [
            {code: 1, name: "端口1"},
            {code: 2, name: "端口2"},
            {code: 3, name: "端口3"},
            {code: 4, name: "端口4"}
        ]
    }
    
    //  版本切换
    $scope.upVersion = {
        title:"版本切换",
        default: {},
        //  选中项
        current: {code: 1, name: "切换为VGA"},
        list: [
            {code: 1, name: "切换为VGA"},
            {code: 2, name: "切换为计量"}
        ]
    }
    //  版本切换 
    $scope.versionTog = {
    	//  版本切换弹窗
    	modalInstance_v:null,
    	title:'版本切换中。。。',
    	numb:0,
    	//  延迟时长 毫秒
    	time:10000
    	//    定时
    }
    $scope.versionTimes = function(item){
		$interval(function(){
    		$scope.versionTog.numb = $scope.versionTog.numb + 1;
    		if($scope.versionTog.numb >= 100){
    			$scope.versionTog.title = '切换完成！';
    			$timeout(function(){
    				$scope.versionTog.modalInstance_v.dismiss();
    				$scope.upVersion.current = item;
    			},2000);
    			
    		}
    	},100,100);
	}
    
    //  版本切换 change
    $scope.fnChangeUpVersion = function(item){
    	$scope.versionTog.title = '版本切换中。。。';
    	$scope.versionTog.numb = 0;
    	//
    	$scope.versionTog.modalInstance_v = $modal.open({
			templateUrl: 'modal/modal_alert_up.html',
			scope: $scope,
			windowClass:'m-modal-alert3',
			keyboard:false,
			backdrop:'static'
		});
		//
		$scope.versionTimes(item);
    }
    
    
    //   设备状态信息数据
    $scope.deviceData = null;


    

    //获取设备状态
    $scope.fnGetDeviceStatus = function () {
        var url = config.zkmb_config.coapServer + "/device/node/state/device";
		var data = {
			UserID: $scope.uid, 		//用户ID(*字符串)   
			DeviceIDs: [$scope.deviceId], 	//设备IDs(*字符串数组)    
		    Params: "" 	//参数(字符串)
		};
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.deviceData = data.Data[0];
            $scope.status($scope.deviceData);
            fnInitPara();//初始化处理
            //console.log($scope.deviceData)

            fnAutoRefresh();
            $.unblockUI();
        }, function (reason) {
        }, function (update) {
        })
    }
    $scope.fnGetDeviceStatus();
    
    
    //  状态
	$scope.DeviceStatus = {
		//   使用时间
		ServiceTime:0,
		//   灯泡寿命
		Status:null
	};
	$scope.status = function(item){
		if(!angular.isObject(item)) return;
		//  使用时间
		$scope.DeviceStatus.ServiceTime = item.UseTimeAfter + item.UseTimeBefore;
		$scope.DeviceStatus.Status = item.DeviceStatus;
	}
    
	//打开设备
	$scope.fnOpenDevice = function() {
		fnStopAutoRefresh();
		$.blockUI({
			message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
		});
		var url = config.zkmb_config.coapServer + "/device/node/control/switch/on/device";
		var data = {
			UserID:$scope.uid,        //用户ID：字符串   
			DeviceID:$scope.deviceId,    //设备ID：字符串   
			Params:""     //参数：字符串(可选)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			fnManualRefresh();
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
		promise.then(function(data) {
			fnManualRefresh();
			// 更新日志等记录
            info_run();
		}, function(reason) {}, function(update) {})
	}

    //切换端口
    $scope.fnChangePort = function (index) {
    	fnStopAutoRefresh();
        $.blockUI({ message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;ִ执行中...</span></div>',timeout:10000 });
        var url = config.zkmb_config.coapServer + "/device/multiplexer/control/vnriver";
        var data = {
            UserID:$scope.uid,
            DeviceID: $scope.deviceId,
		    InPort: $scope.port.list[index].code, 		// 输入端口(*整型)
		    OutPort: 1, 		// 输出端口(*整型)  
		    Params: "" 	//参数(字符串)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.port.current = $scope.port.list[index]
            fnManualRefresh();
            // 更新日志等记录
            info_run();
        }, function (reason) {
        }, function (update) {
        })
    }

//  //读取当前端口
//  $scope.fnReadPort = function (index) {
//      $.blockUI({ message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../project/zkmb/img/zkmb/loading.gif">&nbsp;ִ执行中...</span></div>',timeout:10000 });
//      var url = config.zkmb_config.coapServer + "/device/controlother";
//      var data = {
//          UId:$scope.uid,
//          DeviceId: $scope.deviceId,
//          UseWhoseCmd: "self",
//          CmdCode: "toggle",
//          Para: '"p1":"'+$scope.fnCmdCombind(0xfe,0xfe,0x00,0x31,$scope.port.list[index].code,0x01,0xaa,0xaa)+'"'
//      };
//      var promise = httpService.ajaxPost(url, data);
//      promise.then(function (data) {
//          $scope.port.current = $scope.port.list[index]
//      }, function (reason) {
//      }, function (update) {
//      })
//  }


    //返回
    $scope.fnBack = function(){
        $window.history.back();
    }

    $scope.fnInit=function(){
        $scope.fnGetDeviceStatus();
    }

//  //命令编联
//  $scope.fnCmdCombind=function(b1,b2,b3,b4,b5,b6,b7,b8){
//      return  "hex"+$scope.b2a(b1)+$scope.b2a(b2)+$scope.b2a(b3)+$scope.b2a(b4)+$scope.b2a(b5)+$scope.b2a(b6)+$scope.b2a(b7)+$scope.b2a(b8);
//  }

//  //十六进制整数到十六进制字符串(oxFF > FF)
//  $scope.b2a = function(b){
//      var tmp = b.toString(16);
//      if(tmp.length == 1)
//      {
//          tmp = "0" + tmp;
//      }else{
//          tmp = "" + tmp;
//      }
//      //将第5个字节拼接到hex后面
//      return tmp;
//  }

    //停止自动刷新
    var fnStopAutoRefresh=function(){
        if (!angular.isUndefined(timer)){
            $timeout.cancel(timer);
            timer = 'undefined';
        }
    }

    //手动刷新
    var fnManualRefresh = function(){
        $timeout(function() {
            $scope.fnGetDeviceStatus();
        }, config.zkmb_config.immediatelyRefreshTime);
    }

    //自动刷新
    var fnAutoRefresh = function(){
        timer = $timeout(function() {
            $scope.fnGetDeviceStatus();
        }, config.zkmb_config.fixRefreshTime);
    }

    function fnInitPara(){
        sendContents = $scope.deviceData.LastSendContent;
        for (i=0;i< sendContents.length;i++){
            sendContent = sendContents[i]
            if (sendContent.CmdCode=="toggle" && sendContent.Value!=""){
                $scope.port.current = $scope.port.list[parseInt(sendContent.Value)]
            }
        }

        //从返回的数据中获得VGA设备的当前输入端口
        //a =  $scope.deviceData.DeviceStatus;
        //for (i=0;i< a.length;i++){
        //    if (i.StatusCode=="port"){
        //        //取出port
        //        val = parseInt(i.StatusValueCode)
        //        $scope.port.current = $scope.port.list[val-1]
        //    }
        //}
    }
    
    //   取消定时器
	$scope.$on("$destroy", function() {
		fnStopAutoRefresh();
		//console.log("lamp")
        if (timer) {
            $timeout.cancel(timer);
        }
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
				console.log(data.Reason);
				//alert(data.Reason);
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
				//
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







