/**
 * Created by Administrator on 2016/7/12.
 */
app.controller('devicelistController', function ($scope, httpService,notificService,$location,settings,$interval,$timeout,TipService) {
    $scope.tipService = TipService;
    
    var timer;
    
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
    //user id
    if ($location.search().uid) {
        $scope.uid = $location.search().uid;
    }else{
        $scope.uid = "1";
    }

    //get deviceimg url
    $scope.deviceimg = settings.deviceImg;// "http://192.168.0.102:8080/web/upfile/device/"

    //获取教室信息
    $scope.classroomInfoData ={
      Classroomsname:"",
      Campusname:"",
      Buildingname:"",
      Classroomicon:"",
      Classroomstate:1,
      Curriculumname:"",
      Nickname:""
    };

    $scope.getClassroomInfo = function () {
        var url = settings.getClassroomInfoUrl;
        var data = {
            id: $scope.classroomId
        };
        var promise = httpService.ajaxGet(url, data);
        promise.then(function (data) {
            $scope.classroomInfoData = data.Result;
            console.log(data)
        }, function (reason) {
        }, function (update) {
        })
    }

    //获取设备状态信息
    $scope.fnGetDeviceStatus = function () {
        var url = settings.coapServer + "/device/node/state/room";
        var data = {
            UserID: $scope.uid, 		//用户ID(*字符串)   
            RoomID: $scope.classroomId, 	//房间ID(*字符串)    
            Params: "" 	//参数(字符串)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.deviceData = data.Data;
          	fnAutoRefresh();
            $.unblockUI();
        }, function (reason) {
        }, function (update) {
        })
    }

    //打开
    $scope.fnOpenDevice = function () {
        fnStopAutoRefresh();
        $.blockUI({ message: '<div style="padding: 2px"><span style="font-size: 13px;"> <img src="../img/loading.gif">&nbsp;执行中...</span></div>',timeout:10000 });
        var url = settings.coapServer + "/device/node/control/switch/on/room";
        var data = {
            "UserID": $scope.uid,	//(*字符串)
            "RoomID": $scope.classroomId, //(*字符串)
            "Params": ""	//(字符串)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            fnManualRefresh();
        }, function (reason) {
        }, function (update) {
        })
    }

    //关闭
    $scope.fnCloseDevice = function () {
        fnStopAutoRefresh();
        $.blockUI({ message: '<div style="padding: 2px"><span style="font-size: 13px;"> <img src="../img/loading.gif">&nbsp;执行中...</span></div>',timeout:10000 });
        var url = settings.coapServer + "/device/node/control/switch/off/room";
        var data = {
            "UserID": $scope.uid,	//(*字符串)
            "RoomID": $scope.classroomId, //(*字符串)
            "Params": ""	//(字符串)
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            fnManualRefresh();
        }, function (reason) {
        }, function (update) {
        })
    }

    //打开设备详细页面
    $scope.fnIntoDevice = function (did,page) {
        window.location.href = page +"?cid="+$scope.classroomId+"&did="+did+"&uid="+$scope.uid
    }

    //返回
    $scope.fnBack = function(){
        window.location.href="finishDeviceListPage"
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

    $scope.fnInit=function(){
        $scope.fnGetDeviceStatus();
        $scope.getClassroomInfo();
    }
    
   
   
   
   /*    2016-11-16    */
  //   标题名称
   $scope.Modelname = "设备管理";
   if(config.GetUser().Modelname != ""){
   		$scope.Modelname = config.GetUser().Modelname;
   }
   
   
})







