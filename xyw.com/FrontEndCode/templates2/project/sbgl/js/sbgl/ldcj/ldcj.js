//开启路由器
app.controller("sbglLdcjContr", ['$scope', 'httpService', '$modal', 'toaster',function ($scope, httpService, $modal,toaster) {
	console.log("联动场景");
	$scope.SbglLdcjItems = [];
	//   page
	$scope.backPage = {
		PageIndex:1,
		PageSize:15
	}
	$scope.form = {
		"KeyWord":"",
		"EventSetTableId":0,
		"CampusId":0,
		"BuildingId":0,
		"FloorsId":0,
		"ClassRoomId":0,
		"IsME":0
	}
	//	查询列表
	var getDeviceModelControlCMDList = function(){
		var url = config.HttpUrl + "/Task/TimedTaskList";
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: $scope.backPage.PageIndex,
				PageSize:$scope.backPage.PageSize
			},
			Para:{
				KeyWord:$scope.form.KeyWord,
				CampusId:$scope.form.CampusId,
				BuildingId:$scope.form.BuildingId,
				FloorsId:$scope.form.FloorsId,
				ClassRoomId:$scope.form.ClassRoomId,
				IsME:$scope.form.IsME
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("联动场景查询列表",data);
			if(data.Rcode == "1000"){
				$scope.SbglLdcjItems = data.Result.PageData;
				//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
			}else if(data.Rcode=="1002"){
				$scope.SbglLdcjItems = [];
				//   分页
				var objPage={PageCount:0,PageIndex:1,PageSize:10,RecordCount:0};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
			} else{
				toaster.pop('warning',data.Reason);
			}
		}, function (reason) {}, function (update) {});
	}
	/*  -------------------- 分页、页码  -----------------------  */
	//$scope.backPage = {};
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
	};
	//  翻页
	$scope.pageClick = function(pageindex){
		if(!(Number(pageindex) > 0))return false;
		if(pageindex > 0 && pageindex <= $scope.backPage.PageCount){
			$scope.backPage.PageIndex = pageindex;
			getDeviceModelControlCMDList();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */
	//   查询
	$scope.searchPost = function(){
		getDeviceModelControlCMDList();
	}
	//   回车查询
	$scope.sbgzKeyup = function(e){
		var keycode = window.event?e.keyCode:e.which;
		if(keycode==13){
			getDeviceModelControlCMDList();
		}
	}
	//  删除列表
	var deleteDeviceModelControlCMD = function(TaskId){
		var url = config.HttpUrl + "/Task/DelTimedTask";
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para:{
				TaskId:TaskId
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("删除成功",data);
			if(data.Rcode =="1000"){
				getDeviceModelControlCMDList();
        toaster.pop('success',"删除成功！");
			}else{
        toaster.pop('warning',data.Reason);
			}
		});
	}
	/*删除*/
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'确定删除设备管理联动场景？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        deleteDeviceModelControlCMD(item.TaskId);
      }
    });
  }
	//  添加按钮功能
	$scope.openModalAddMl = function (item,str) {
		if(!str)str = "";
		if(!item)item = {};
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/ldcj/modal_ldcj.html',
			controller: 'modalSbglLdcjContr',
			windowClass: 'm-modal-sbgl-sbpz',
			resolve: {
				items: function () {
					return {"operate":str,"item":item};
				}
			}
		});
		modalInstance.result.then(function(bul){
			console.log("bul",bul);
			if(bul){
				getDeviceModelControlCMDList();
			}
		});
	}

	$scope.run = function(){
		getDeviceModelControlCMDList();
	}
	$scope.run();
}]);

//设备管理-联动场景-添加联动场景弹窗
app.controller("modalSbglLdcjContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("设备管理-联动场景-添加联动场景弹窗");
	$scope.items = items;
	console.log("弹窗",$scope.items);
	$scope.form = {
		"Campusname":"",
		//	定时任务ID
		"TaskId":0,
    	//  打开/关闭定时任务
		"IsOpen":1,
		//  -----   打开/关闭定时任务
		"IsOpenItem":'',
		//  定时任务开关
		"IsOpenItems": [
			{"val":0,"title":"关闭"},
			{"val":1,"title":"打开"}
		],
		//	定时任务名称
		"TaskName":"",
	    //  定时任务开关ID 0 关闭 1 打开
	    "TaskIsOpen":0,
		//	选择的时间点
		"TimePoint":"",
		//	重复类型
		"RepeatType":"",
		//  自定义选择的值
		"RepeatValue":"",
		//	自定义触发条件显示名称
		"repeatText":"",
		//  响应的事件Id
		"EventSetTableId":0,
		//	教室Id
		"ClassRoomId":0,
		//	楼栋Id
		"BuildingId":0,
		//	楼层Id
		"FloorsId":0,
		//	校区Id
		"CampusId":0,
		//	是否查询到楼层
		"IsFloors":"0",
    //	是否查询到教室
		"IsClassRoom":"0",
    //  节点id[保留扩展]
		"NodeId":"",
    //  设备id[保留扩展]
		"DeviceId":"",
    //  命令id[保留扩展]
		"CmdId":"",
		//	位置名称
		"postionName":"",
    //  位置编号
    	"postionCode":""
	}


	$scope.changeIsOpenItem = function(item){
		$scope.form.IsOpenItem = item;
		$scope.form.TaskIsOpen = item.val;
	}

	$scope.XysjItem = "";
	$scope.XysjItems = [];
	//  ********查询列表*******************
	$scope.changeXysjItem = function(item){
		$scope.form.EventSetTableId = item.EventSetTableId;
	}
	
	var EventSetTableList = function(){
    var url = config.HttpUrl+"/Task/EventSetTableList";
		var data = {
      Auth:{
        "Usersid": config.GetUser().Usersid,
        "Rolestype": config.GetUser().Rolestype,
        "Token": config.GetUser().Token,
        "Os": "WEB"
      },
      Page:{
        "PageIndex":1,
        "PageSize":100
      },
      Para:{
        "IsFloors":$scope.form.IsFloors,
        "IsClassRoom":$scope.form.IsClassRoom,
        "FloorsId":Number($scope.form.FloorsId),
        "ClassRoomId":Number($scope.form.ClassRoomId),
        "NodeId":$scope.form.NodeId,
        "DeviceId":$scope.form.DeviceId,
        "CmdId":$scope.form.CmdId
      }
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("查询响应事件列表",data);
			if(data.Rcode =="1000"){
				$scope.XysjItems = data.Result.PageData;
				for(var item in $scope.XysjItems) {
					if($scope.XysjItems[item].EventSetTableId == $scope.form.EventSetTableId) {
						$scope.XysjItem = $scope.XysjItems[item];
					}
				}
			}else{
				toaster.pop('warning',data.Reason);
			}
		});
	}
  //	查询打开/关闭定时任务
  var getDeviceModelControlOnList = function(){
    var url = config.HttpUrl + "/Task/OnOrOffTimedTask";
    var data = {
      Auth:{
        "Usersid": config.GetUser().Usersid,
        "Rolestype": config.GetUser().Rolestype,
        "Token": config.GetUser().Token,
        "Os": "WEB"
      },
      Para:{
        TaskId:Number($scope.form.TaskId),
        IsOpen:Number($scope.form.IsOpen)
      }
    }
    var promise = httpService.ajaxPost(url,data);
    promise.then(function(data){
    
    }, function (reason) {}, function (update) {});
  }
	//联动场景--添加
	var saveDeviceModelControlCMDAdd = function(){
		if(!(formValidate($scope.form.TaskName).minLength(0).outMsg(2700).isOk))return false;
		if(!(formValidate($scope.form.repeatText).minLength(0).outMsg(2701).isOk))return false;
		if(!(formValidate($scope.form.EventSetTableId).isNumber(0).outMsg(2702).isOk))return false;
		if(!(formValidate($scope.form.postionName).minLength(0).outMsg(2703).isOk))return false;

		var url = config.HttpUrl+"/Task/AddTimedTask";
		var data={
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				"TaskName": $scope.form.TaskName,
				"TimePoint":$scope.form.TimePoint,
				"RepeatType":$scope.form.RepeatType,
				"RepeatValue":$scope.form.RepeatValue,
				"EventSetTableId":Number($scope.form.EventSetTableId),
				"ClassRoomId":Number($scope.form.ClassRoomId),
				"BuildingId":Number($scope.form.BuildingId),
				"FloorsId":Number($scope.form.FloorsId),
				"CampusId":Number($scope.form.CampusId),
				"TaskIsOpen":Number($scope.form.TaskIsOpen)
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("联动场景-添加",data);
      if(data.Rcode == "1000"){
        toaster.pop('success', '添加成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		},function(reason){},function(update){});
	}
	//  添加按钮功能
	$scope.openModalAddMlAddTJ = function (item,str) {
		if(!str)str = "";
		if(!item)item = {};
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/ldcj/modal_ldcj_add.html',
			controller: 'modalSbglLdcjAddContr',
			windowClass: 'm-modal-sbgl-sbpz',
			resolve: {
				items: function () {
					return {"operate":str,"item":item};
				}
			}
		});
		modalInstance.result.then(function(obj){
			console.log("obj",obj);
			if(obj){
				$scope.form.repeatText = "";
				$scope.form.TimePoint = obj.UploadTime;
				$scope.form.RepeatValue = obj.weekValue;
				$scope.form.RepeatType= obj.rpted;
				if(obj.UploadTime){$scope.form.repeatText += "（" + obj.UploadTime + '）';}
				if(obj.rpted){$scope.form.repeatText += "（" + obj.rpted + '）';}
				if(obj.weekValue){$scope.form.repeatText += "（" + obj.weekValue + '）';}
			}
		});
	}
	//控制命令--修改
	var saveDeviceModelControlCMD = function(){
		if(!(formValidate($scope.form.TaskName).minLength(0).outMsg(2700).isOk))return false;
		if(!(formValidate($scope.form.repeatText).minLength(0).outMsg(2701).isOk))return false;
		if(!(formValidate($scope.form.EventSetTableId).isNumber(0).outMsg(2702).isOk))return false;
		if(!(formValidate($scope.form.postionName).minLength(0).outMsg(2703).isOk))return false;
		var url = config.HttpUrl+"/Task/ChangeTimedTask";
		var data={
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				"TaskId": Number($scope.form.TaskId),
				"TaskName":$scope.form.TaskName,
				"TimePoint":$scope.form.TimePoint,
				"repeatText":$scope.form.repeatText,
				"RepeatType":$scope.form.RepeatType,
				"RepeatValue":$scope.form.RepeatValue,
				"EventSetTableId":Number($scope.form.EventSetTableId),
				"ClassRoomId":Number($scope.form.ClassRoomId),
				"BuildingId":Number($scope.form.BuildingId),
				"FloorsId":Number($scope.form.FloorsId),
				"CampusId":Number($scope.form.CampusId),
        		"TaskIsOpen":Number($scope.form.TaskIsOpen)
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("联动场景-修改",data);
	      if(data.Rcode == "1000"){
	        toaster.pop('success', '修改成功！');
	        $modalInstance.close(true);
	      }else{
	        toaster.pop('warning',data.Reason);
	      }
		},function(reason){},function(update){});
	}
	//保存按钮
	$scope.ok = function(){
		if($scope.items.operate == "edit"){
			saveDeviceModelControlCMD();
      getDeviceModelControlOnList();
		}else{
			saveDeviceModelControlCMDAdd();
		}
	}
	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}
	//    打开弹窗  选择学校
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
      console.log(selectedItem);

      if(selectedItem.addCode == 'floor' || selectedItem.addCode =='classroom'){
              var items = selectedItem.addItems;
              if( "campus" in items){
	              $scope.form.CampusId = items.campus.addId;
	            }else {
	              $scope.form.CampusId = 0;
	            }
	            if( "building" in items){
	              $scope.form.BuildingId = items.building.addId;
	            }else {
	              $scope.form.BuildingId = 0;
	            }
              
              $scope.form.IsFloors = "0";
              $scope.form.IsClassRoom = "0";
              if(selectedItem.addCode == 'floor'){
              	$scope.form.FloorsId = selectedItem.addId;
              	$scope.form.ClassRoomId = 0;
              }else{
              	$scope.form.FloorsId = 0;
              	$scope.form.ClassRoomId = selectedItem.addId;
              }
              
      	
          $scope.form.postionName = selectedItem.add;
          $scope.form.postionCode = selectedItem.addCode;
          
          EventSetTableList();
          
       
	      }else if(selectedItem.addCode == 'campus' || selectedItem.addCode =='building' || selectedItem.addCode == ''){
          $modal.open({
            templateUrl: 'modal/modal_alert_all.html',
            controller: 'modalAlert2Conter',
            resolve: {
              items: function () {
                return {"type":'warning',"msg":'至少选择到楼层!'};
              }
            }
          });
      }
      

		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	$scope.run = function(){
		for(var item in $scope.form.IsOpenItems) {
			if($scope.form.IsOpenItems[item].val == $scope.items.item.TaskIsOpen) {
				$scope.form.IsOpenItem = $scope.form.IsOpenItems[item];
			}
		}
		if(items.operate == "edit" || items.operate == "details"){
			$scope.form = $.extend({},$scope.form,$scope.items.item);
		      if( $scope.items.item.Classroomsname.length > 0){
				  $scope.form.postionName = $scope.items.item.Campusname + '-' + $scope.items.item.BuildingName + '-' + $scope.items.item.Floorname +'-'+ $scope.items.item.Classroomsname;
		      }else{
				  $scope.form.postionName = $scope.items.item.Campusname + '-' + $scope.items.item.BuildingName + '-' + $scope.items.item.Floorname;
		      }
			$scope.form.EventSetTableId = $scope.items.item.EventSetTableId;
			$scope.form.TaskId = $scope.items.item.TaskId;
			$scope.form.IsOpen = $scope.items.item.TaskIsOpen;
			if($scope.items.item.RepeatValue !== ""){
				$scope.form.repeatText = "（" + $scope.items.item.TimePoint + "）（" + $scope.items.item.RepeatType + "）（" + $scope.items.item.RepeatValue + "）";
			}else{
				$scope.form.repeatText = "（" + $scope.items.item.TimePoint + "）（" + $scope.items.item.RepeatType + "）";
			}
		}
		EventSetTableList();
	}
	$scope.run();
}]);
//设备管理-联动场景-添联动场景弹窗-触发条件弹窗
app.controller("modalSbglLdcjAddContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate',function ($scope, httpService, $modal,$modalInstance,items,formValidate) {
	console.log("设备管理-联动场景-添联动场景弹窗-触发条件弹窗");
	$scope.items = items;
	console.log("触发条件弹窗",$scope.items);
	$scope.form = {
		"rpted":"",
		"rptedValue":"",
		//    重复选中的对象
		"rptedValueItem":"",
		"weekValue":"",
		"UploadTime":""
	}
	//触发条件
	$scope.timing = [
		{value:0, name:"只执行一次"},
		{value:1, name:"每天"},
		{value:2, name:"自定义"}
	];
	//触发条件-自定义周期
	$scope.Aweek = [
		{value:0, name:"星期一"},
		{value:1, name:"星期二"},
		{value:2, name:"星期三"},
		{value:3, name:"星期四"},
		{value:4, name:"星期五"},
		{value:5, name:"星期六"},
		{value:6, name:"星期天"}
	]
	//	取消按钮
	$scope.cancel = function(){
		$modalInstance.dismiss('cancel');
	}
	//    打开弹出 -选择日期
	$scope.showDate = function() {
		jeDate({
			dateCell: "#jd_begindate",
			format: "hh:mm",
			isTime: true,
			minDate: "00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.UploadTime = val;
			},
			okfun: function(elem,val) {
				$scope.form.UploadTime = val;
			},
			clearfun:function(elem, val) {
				$scope.form.UploadTime = "";
			}
		});
	}
	//    重复
	$scope.changeSelect = function(item) {
		$scope.form.rptedValueItem = item;
		$scope.form.rptedValue = item.value;
		//
		$scope.form.modalOpenClassroom = item.name;
		$scope.form.rpted = item.name;
	}
	$scope.ok = function() {
		$scope.form.weekValue = "";
		angular.forEach($scope.Aweek,function(n,v){
			if('checked' in n && n.checked == true){
				$scope.form.weekValue += n.name + "|";
			}
		});
		if($scope.form.weekValue != ""){
			$scope.form.weekValue = $scope.form.weekValue.substr(0,$scope.form.weekValue.length-1);
		}
		if($scope.form.rpted != "自定义"){
			$scope.form.weekValue = "";
		}
		$modalInstance.close($scope.form);
	}
	if(items.operate == "edit" || items.operate == "details"){
		$scope.form.UploadTime = $scope.items.item.TimePoint;
		$scope.form.rpted = $scope.items.item.RepeatType;
		$scope.form.weekValue = $scope.items.item.RepeatValue;
		var c = $scope.form.weekValue.split("|");
		for(var i in $scope.Aweek){
			//	循环拆分插入选中项
			for(var a in c){
				if($scope.Aweek[i].name == c[a]){
					$scope.Aweek[i].checked = true;
				}
			}
		}
	}
}]);
