/**
 * Created by Administrator on 2016/7/29.
 */
app.controller('vgatoggleController', function ($scope, httpService, notificService, $location, settings, $window, $interval, $timeout, $filter, $modal, TipService) {

    $scope.port = {
        title:"切换端口",
        default: {},
        current: {},
        list: [
            {code: 0x01, name: "端口1"},
            {code: 0x02, name: "端口2"},
            {code: 0x03, name: "端口3"},
            {code: 0x04, name: "端口4"}
        ]
    }
	
	//   是否有NAV颜色
    if($location.search().navcolor){
    	$scope.navcolor = '#' + $location.search().navcolor;
    }

    //classroom id
    if ($location.search().cid) {
        $scope.classroomId = $location.search().cid;
    }else{
        $scope.classroomId = "1";
    }
    //device id
    if ($location.search().did) {
        $scope.deviceId = $location.search().did;
    }else{
        $scope.deviceId = "1";
    }
    //user id
    if ($location.search().uid) {
        $scope.uid = $location.search().uid;
    }else{
        $scope.uid = "0";
    }
    //  oney 是否为单页， 用于移动端 ，单页可直接返回
    if ($location.search().uid) {
        $scope.only = $location.search().only;
    }else{
        $scope.only = "0";
    }

    //获取设备状态
    $scope.fnGetDeviceStatus = function () {
        var url = settings.coapServer + "/device/node/state/device";
        var data = {
            UserID: $scope.uid, 		//用户ID(*字符串)   
            DeviceIDs: [$scope.deviceId], 	//设备IDs(*字符串数组)    
            Params: "" 	//参数(字符串)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.deviceData = data.Data[0];

            fnInitPara();//初始化处理

            fnAutoRefresh();
            $.unblockUI();
        }, function (reason) {
        }, function (update) {
        })
    }


    //打开设备
    $scope.fnOpenDevice = function () {
        fnStopAutoRefresh();
        $.blockUI({
            message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../img/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
        });
        var url = settings.coapServer + "/device/node/control/switch/on/device";
        var data = {
            UserID: $scope.uid,        //用户ID：字符串   
            DeviceID: $scope.deviceId,    //设备ID：字符串   
            Params: ""     //参数：字符串(可选)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            fnManualRefresh();
            // 更新日志等记录
            info_run();
        }, function (reason) { }, function (update) { })
    }

    //关闭
    $scope.fnCloseDevice = function () {
        fnStopAutoRefresh();
        $.blockUI({
            message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../img/loading.gif">&nbsp;执行中...</span></div>',timeout:10000
        });
        var url = settings.coapServer + "/device/node/control/switch/off/device";
        var data = {
            UserID: $scope.uid,        //用户ID：字符串   
            DeviceID: $scope.deviceId,    //设备ID：字符串   
            Params: ""     //参数：字符串(可选)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            fnManualRefresh();
            // 更新日志等记录
            info_run();
        }, function (reason) { }, function (update) { })
    }



    //切换端口
    $scope.fnChangePort = function (index) {
        $.blockUI({ message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../img/loading.gif">&nbsp;ִ执行中...</span></div>',timeout:10000 });
        var url = config.zkmb_config.coapServer + "/device/multiplexer/control/vnriver";
        var data = {
        	UserID: $scope.uid, 		//用户ID(*字符串)   
		    DeviceID: $scope.deviceId, 	//设备ID(*字符串) 
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

    //读取当前端口
    $scope.fnReadPort = function (index) {
        $.blockUI({ message: '<div style="padding: 10px"><span style="font-size: 13px;"> <img src="../img/loading.gif">&nbsp;ִ执行中...</span></div>',timeout:10000 });
        var url = settings.coapServer + "/device/controlother";
        var data = {
            UId:$scope.uid,
            DeviceId: $scope.deviceId,
            UseWhoseCmd: "self",
            CmdCode: "toggle",
            Para: '"p1":"'+$scope.fnCmdCombind(0xfe,0xfe,0x00,0x31,$scope.port.list[index].code,0x01,0xaa,0xaa)+'"'
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.port.current = $scope.port.list[index]
        }, function (reason) {
        }, function (update) {
        })
    }


    //返回
    $scope.fnBack = function(){
        if($scope.only == "1"){
    		window.location.href="finishDeviceListPage";
    	}else{
       		$window.history.back();
    	}
    }

    $scope.fnInit=function(){
        $scope.fnGetDeviceStatus();
    }

    //命令编联
    $scope.fnCmdCombind=function(b1,b2,b3,b4,b5,b6,b7,b8){
        return  "hex"+$scope.b2a(b1)+$scope.b2a(b2)+$scope.b2a(b3)+$scope.b2a(b4)+$scope.b2a(b5)+$scope.b2a(b6)+$scope.b2a(b7)+$scope.b2a(b8);
    }

    //十六进制整数到十六进制字符串(oxFF > FF)
    $scope.b2a = function(b){
        var tmp = b.toString(16);
        if(tmp.length == 1)
        {
            tmp = "0" + tmp;
        }else{
            tmp = "" + tmp;
        }
        //将第5个字节拼接到hex后面
        return tmp;
    }

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
        }, settings.immediatelyRefreshTime);
    }

    //自动刷新
    var fnAutoRefresh = function(){
        timer = $timeout(function() {
            $scope.fnGetDeviceStatus();
        }, settings.fixRefreshTime);
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
    
    
    
    
    /*  ---------------------   new 2016-10-18  ---------------------------------  */
    $scope.tipService = TipService;
	//    图片地址
	$scope.deviceimg = settings.deviceImg;
	//    教室数据
	$scope.classRoomItem = {};
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
	
	
	//    取教室信息
	var getClassroomInfo = function(classroomid) {
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
			    $scope.tipService.setMessage(data.Reason);
				//alert(data.data);
			}
		}, function(reason) {}, function(update) {});
	}
	
	
	/*--------------------------  设备数据   ------------------------------*/
	//   默认第一页
	var PageIndex = 1;
	//   分页页显示条数
	var PageSize = 10;
	//   显示页码数
	var pageNumber = 3;
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
				//console.log("预警")
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
			    $scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("设备故障")
			//console.log(data)
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
			    $scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("设备操作日志")
			//console.log(data)
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
			    $scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("设备使用日志")
			//console.log(data)
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
		    if (data.Rcode == "1000") {
		        $scope.tipService.setMessage("删除成功！");
				//alert("删除成功！");
				getFaultInfoList($scope.backPage3.PageIndex);
			} else {
			    $scope.tipService.setMessage(data.Reason);
				//alert(data.Reason);
			}
			//console.log("删除故障")
			//console.log(data)
		}, function(reason) {}, function(update) {});
	}
	
	//   故障添加打开弹窗
	$scope.modalOpenAddFault = function(item,str) {
		//    当登录用户不为故障登记用户时不能修改删除
		if(str == 'delete' || str == 'edit'){
			if(item.InputUserId != config.GetUser().Usersid && item.Status == 0){
				$scope.tipService.setMessage("登录账号错误！");
				return void(0);
			}
		}
		
		
		//
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
				}
				/*,
				deps: ['$ocLazyLoad',
					function($ocLazyLoad) {
						return $ocLazyLoad.load(['ui.select']);
					}
				]*/
			}
		});
		
		modalInstance.result.then(function(bul) {
			//console.log(bul)
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
		getClassroomInfo($scope.classroomId);
		getAlertInfoList();
		getFaultInfoList(PageIndex);
		getOperateLogList(PageIndex);
		getUseLogList(PageIndex);
	}
	info_run();
	//  
	
	//   分页
	
	/*--------------------------  设备数据 End   ------------------------------*/
	
	
	
    /*  ---------------------   new End   ---------------------------------  */
    
    
    

});







