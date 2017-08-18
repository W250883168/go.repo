//开启路由器
app.controller("sbglSbglzContr", ['$scope', 'httpService', '$modal','toaster', function ($scope, httpService, $modal,toaster) {
	console.log("设备配置-设备管理者");

	$scope.form = {
		//    查询关键字:字符串非必填
    	"KeyWord":'',
		//    设备型号ID:字符串 非必填
        "ModelId":"",
		//    设备型号名称
        "ModelName":"",
		//    楼栋ID:字符串 非必填
        "Buildingid":0,
		//    楼层ID:字符串 非必填
        "Floorsid":0,
		//    校区ID:字符串 非必填
        "Campusid":0,
		//    教室ID:字符串 非必填
        "ClassroomId":0,
		//    安装位置
        "addHtml":""
	}

	//   设备列表
	$scope.deviceList = [];

	//    分页
	$scope.backPage = {
        PageIndex:1,
        PageSize:10
	}

	/**
	 * 取设备列表
	 */
    var getDeviceList = function(){
		var url = config.HttpUrl + "/device/getDeviceList";
		var data = {
            "Auth":{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
            Page: {
				PageIndex:$scope.backPage.PageIndex,
				PageSize:$scope.backPage.PageSize
			},
            "Para":{
                //    查询关键字:字符串非必填
		    	"KeyWord":$scope.form.KeyWord,
		    	//    设备型号ID:字符串 非必填
		        "ModelId":$scope.form.ModelId,
		        //    楼栋ID:字符串 非必填
		        "Buildingid":Number($scope.form.Buildingid),
		        //    楼层ID:字符串 非必填
		        "Floorsid":Number($scope.form.Floorsid),
		        //    校区ID:字符串 非必填
		        "Campusid":Number($scope.form.Campusid),
		        //    教室ID:字符串 非必填
		        "ClassroomId":Number($scope.form.ClassroomId)
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("取设备列表SS",data);
            if(data.Rcode =="1000"){
                $scope.deviceList = data.Result.Data;
                //分页
                $scope.backPage = pageFn(data.Result.Page,5);
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }


	/**
	 * 删除设备
	 */
    var deleteDevice = function(Id){
		var url = config.HttpUrl + "/device/deleteDevice";
		var data = {
            "Auth":{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
            "Para":{
				//    查询关键字:字符串非必填
		    	"Id":Id
			}
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("删除设备",data);
            if(data.Rcode =="1000"){
				getDeviceList();
				toaster.pop('success', '删除成功！');
            }else{
              toaster.pop('warning',data.Reason);
			}
		});
	};


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
    $scope.pageClick = function(pageindex){
			Number(pageindex) < 1 ? pageindex = 1 : pageindex = Number(pageindex);
			$scope.backPage.PageIndex = pageindex;
			getDeviceList();
    }
		/*  -------------------- 分页、页码  -----------------------  */



	//添加按钮功能
    $scope.openModalAddGlz = function (str,item) {
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbpz/modal_sbglz.html',
			controller: 'modalSbglzContr',
			windowClass: 'm-modal-sbgl-sbpz',
			resolve: {
                items: function () {
                    return {"operate":str,"item":item};
				}
			}
		});

		//    弹窗返回
		modalInstance.result.then(function(bol) {
			console.log(bol)
			if(!bol){
				//
			}else{
				//    刷新列表
				getDeviceList();
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});

	};

	//    打开弹窗  选择设备型号
	$scope.modalOpenDevice = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_device.html',
			controller: 'modalGetDeviceCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(deviceItem) {
			console.log(deviceItem)
			if(!deviceItem){
				$scope.form.ModelId = "";
				$scope.form.ModelName = "";
			}else{
				$scope.form.ModelId = deviceItem.Id;
				$scope.form.ModelName = deviceItem.Name;
			}
			//
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择教室
    $scope.modalOpenClassroom =function(){
		console.log("打开弹窗 -选择教室");
    	var modalInstance=$modal.open({
    		templateUrl:'../html/modal/modal_school.html',
    		controller:'modalGetClassRoomCtrl',
    		resolve:{
    			items:function(){
					return $scope.itens;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log("qq",selectedItem)
			if(selectedItem.addId == ""){
				$scope.form.Campusid = 0;
				$scope.form.Buildingid = 0;
				$scope.form.Floorsid = 0;
				$scope.form.ClassroomId = 0;
				$scope.form.addHtml = "";
			}else{
				switch(selectedItem.addCode){
					case "campus":
						//
						$scope.form.Campusid = selectedItem.addId;
						$scope.form.Buildingid = 0;
						$scope.form.Floorsid = 0;
						$scope.form.ClassroomId = 0;
						$scope.form.addHtml = selectedItem.add;
						break;
					case "building":
						//
						$scope.form.Campusid = 0;
						$scope.form.Buildingid = selectedItem.addId;
						$scope.form.Floorsid = 0;
						$scope.form.ClassroomId = 0;
						$scope.form.addHtml = selectedItem.add;
						break;
					case "floor":
						//
						$scope.form.Campusid = 0;
						$scope.form.Buildingid = 0;
						$scope.form.Floorsid = selectedItem.addId;
						$scope.form.ClassroomId = 0;
						$scope.form.addHtml = selectedItem.add;
						break;
					case "classroom":
						//
						$scope.form.Campusid = 0;
						$scope.form.Buildingid = 0;
						$scope.form.Floorsid = 0;
						$scope.form.ClassroomId = selectedItem.addId;
						$scope.form.addHtml = selectedItem.add;
						break;
				}
			}
			//
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    查询
    $scope.searchPost = function(){
			$scope.backPage.PageIndex = 1;
			getDeviceList();
		}
		//   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
				getDeviceList();
			}
		}
		//  删除弹窗功能
  $scope.deleteItem = function(item){
		var modalInstance = $modal.open({
			templateUrl: 'modal/modal_alert_all.html',
			controller: 'modalAlert2Conter',
			resolve: {
        items: function () {
          return {"type":'warning',"msg":'您确定要删除该设备吗？'};
				}
			}
		});
    modalInstance.result.then(function(bul){
      if(bul){
								deleteDevice(item.Id);
							}
						});
					}
    $scope.run = function(){
		getDeviceList();
	}
	$scope.run();

}]);

//设备管理-设备配置-节点型号-添加节点型号弹窗
app.controller("modalSbglzContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("设备管理-设备配置-设备管理者-添加设备管理者弹窗",items);

	$scope.items = items;

	$scope.form = {
		//    设备Id 必填，前端js参数guid
		"Id": "",
		//    设备名称 必填
		"Name": "",
		//    设备序列号 必填
		"Sn": "",
		//    设备编码 必填
		"Code": "",
		//    设备品牌 必填
		"Brand": "",
		//    设备型号Id 必填
		"ModelId": "",
		//    设备型号名称HTML
		"ModelName": "",
		//    所在教室Id 必填
		"ClassroomId": null,
		//    所在教室名称HTML
		"ClassroomName": "",
		//    电源节点Id 非必填
		"PowerNodeId": "",
		//    设备连接节点开关Id 非必填
        "PowerSwitchId": "1",
		//    接入方式 非必填
		"JoinMethod": "node",
		//  ----  接入方式 非必填
		"JoinMethodItem": "",
		//  ----  接入方式 非必填
		"JoinMethodItems": [],
		//    接入节点Id 非必填
		"JoinNodeId": "",
		//    节点开关状态 非必填
        "JoinSocketId": "1",
		//    非必填
		"NodeSwitchStatus": "",
		//  ----  节点开关状态
		"NodeSwitchStatusItem": "",
		//  ----  节点开关状态
		"NodeSwitchStatusItems": [],
		//    自身开关状态 非必填
		"DeviceSelfStatus": "1",
		//    设备是否可用 非必填[1:可用，0:不可用]
		"IsCanUse": "1",
		//    录入系统前设备已使用时间（秒） 非必填
		"UseTimeBefore": null,
		//    录入系统后设备已使用时间 非必填
		"UseTimeAfter": null
	}
	
	//   选择节点变量
	$scope.nodeIds = {
		item:"",
		items:[]
	}


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

	//    ----  节点开关状态
	$scope.changeNodeSwitchStatusItem = function(item){
		$scope.form.NodeSwitchStatusItem = item;
		$scope.form.NodeSwitchStatus = item.val;
	}

	//    ----  接入方式
	$scope.changeJoinMethodItem = function(item){
		$scope.form.JoinMethodItem = item;
		$scope.form.JoinMethod = item.val;
	}
	
	//   ----  change 选择节点
	$scope.changeNodeItem = function(item){
		$scope.nodeIds.item = item;
		//
		$scope.form.PowerNodeId = item.Id;
		$scope.form.JoinNodeId = item.Id;
	}
	
	//   change 电源节点输入发生变化   同步到 通讯节点 id
	$scope.changePowerNodeId = function(id){
		$scope.form.JoinNodeId = id;
	}
	
	
	/**
	 * 查找 获取节点列表
	 * 
	 */
	var getNodeList = function(){

        var url = config.HttpUrl + "/device/getNodeList";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            Page: {
				PageIndex:-1
			},
            "Para":{
                KeyWord: '', 			//查询关键字
				NodeId: '', 			//节点型号ID
				ClassRoomIds: $scope.form.ClassroomId.toString(), 		//教室ID [可多选]
				Floorsids: '',			//楼层ID [可多选]
				Buildingids: '', 		//楼栋IDs[可多选]
				Campusids: '',			//校区ID:字符串[可多选]
				IsNoSave: ''			//是否安装[enum=0:未选/1:已选]
            }
        }
        //console.log(data)
        //return false;
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("添加设备",data);
            if(data.Rcode == "1000"){
              $scope.nodeIds.items = data.Result.Data;
              if($scope.nodeIds.items && $scope.nodeIds.items.length > 0){
              	$scope.nodeIds.item = $scope.nodeIds.items[0];
              }
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }


	/**
	 * 添加设备
	 */
	var saveDevice = function(){
		//   表单验证
		if(!(formValidate($scope.form.Name).minLength(0).outMsg(2516).isOk))return false;
		//if(!(formValidate($scope.form.Code).minLength(0).outMsg(2517).isOk))return false;
		//if(!(formValidate($scope.form.Brand).minLength(0).outMsg(2518).isOk))return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2500).isOk))return false;
		if(!(formValidate($scope.form.ClassroomName).minLength(0).outMsg(2519).isOk))return false;
		//   是否有填一项
		if(!(formValidate($scope.form.PowerNodeId).minLength(0).isOk) && !(formValidate($scope.form.JoinNodeId).minLength(0).isOk)){
			//alert('节点编号或接入方式节点编号必填一项!');
			return false;
		}
		//   当两项都填时必须一至
		if($scope.form.PowerNodeId.length > 0 && $scope.form.JoinNodeId.length > 0){
			if($scope.form.PowerNodeId != $scope.form.JoinNodeId){
				//alert('节点编号与接入方式节点编号必须相同!');
				return false;
			}
		}


		var url = config.HttpUrl + "/device/saveDevice";
		var data = {
            "Auth":{
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
            "Para":{
					//    设备Id 必填，前端js参数guid
					"Id": getGUIDs().toString(),
					//    设备名称 必填
					"Name": $scope.form.Name,
					//    设备序列号 必填
					"Sn": $scope.form.Sn,
					//    设备编码 必填
					"Code": $scope.form.Code,
					//    设备品牌 必填
					"Brand": $scope.form.Brand,
					//    设备型号Id 必填
					"ModelId": $scope.form.ModelId,
					//    所在教室Id 必填
					"ClassroomId": Number($scope.form.ClassroomId),
					//    电源节点Id 非必填
					"PowerNodeId": $scope.form.PowerNodeId,
					//    设备连接节点开关Id 非必填
					"PowerSwitchId": $scope.form.PowerSwitchId,
					//    接入方式 非必填
					"JoinMethod": $scope.form.JoinMethod,
					//    接入节点Id 非必填
					"JoinNodeId": $scope.form.JoinNodeId,
					//    节点开关状态 非必填
					"JoinSocketId": $scope.form.JoinSocketId,
					//    非必填
					"NodeSwitchStatus": $scope.form.NodeSwitchStatus,
					//    自身开关状态 非必填
					"DeviceSelfStatus": $scope.form.DeviceSelfStatus,
					//    设备是否可用 非必填[1:可用，0:不可用]
					"IsCanUse": $scope.form.IsCanUse,
					//    录入系统前设备已使用时间（秒） 非必填
					"UseTimeBefore": Number($scope.form.UseTimeBefore),
					//    录入系统后设备已使用时间 非必填
					"UseTimeAfter": Number($scope.form.UseTimeAfter)
				}
			}
			//console.log(data)
			//return false;
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("添加设备",data);
            if(data.Rcode == "1000"){
				toaster.pop('success', '添加成功！');
				$modalInstance.close(true);
            }else{
              toaster.pop('warning',data.Reason);
			}
		});
	}


	/**
	 * 修改设备
	 */
	var saveDeviceEdit = function(){
		//   表单验证
		if(!(formValidate($scope.form.Name).minLength(0).outMsg(2516).isOk))return false;
		//if(!(formValidate($scope.form.Code).minLength(0).outMsg(2517).isOk))return false;
		//if(!(formValidate($scope.form.Brand).minLength(0).outMsg(2518).isOk))return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2500).isOk))return false;
		if(!(formValidate($scope.form.ClassroomName).minLength(0).outMsg(2519).isOk))return false;

		//   是否有填一项
		if(!(formValidate($scope.form.PowerNodeId).minLength(0).isOk) && !(formValidate($scope.form.JoinNodeId).minLength(0).isOk)){
			//alert('节点编号或接入方式节点编号必填一项!');
			return false;
		}
		//   当两项都填时必须一至
		if($scope.form.PowerNodeId.length > 0 && $scope.form.JoinNodeId.length > 0){
			if($scope.form.PowerNodeId != $scope.form.JoinNodeId){
				//alert('节点编号与接入方式节点编号必须相同!');
				return false;
			}
		}

		var url = config.HttpUrl + "/device/saveDevice";
		var data = {
            "Auth":{
					"Usersid": config.GetUser().Usersid,
					"Rolestype": config.GetUser().Rolestype,
					"Token": config.GetUser().Token,
					"Os": "WEB"
				},
            "Para":{
					//    设备Id 必填，前端js参数guid
					"Id": $scope.form.Id,
					//    设备名称 必填
					"Name": $scope.form.Name,
					//    设备序列号 必填
					"Sn": $scope.form.Sn,
					//    设备编码 必填
					"Code": $scope.form.Code,
					//    设备品牌 必填
					"Brand": $scope.form.Brand,
					//    设备型号Id 必填
					"ModelId": $scope.form.ModelId,
					//    所在教室Id 必填
					"ClassroomId": Number($scope.form.ClassroomId),
					//    电源节点Id 非必填
					"PowerNodeId": $scope.form.PowerNodeId,
					//    设备连接节点开关Id 非必填
					"PowerSwitchId": $scope.form.PowerSwitchId,
					//    接入方式 非必填
					"JoinMethod": $scope.form.JoinMethod,
					//    接入节点Id 非必填
					"JoinNodeId": $scope.form.JoinNodeId,
					//    节点开关状态 非必填
					"JoinSocketId": $scope.form.JoinSocketId,
					//    非必填
					"NodeSwitchStatus": $scope.form.NodeSwitchStatus,
					//    自身开关状态 非必填
					"DeviceSelfStatus": $scope.form.DeviceSelfStatus,
					//    设备是否可用 非必填[1:可用，0:不可用]
					"IsCanUse": $scope.form.IsCanUse,
					//    录入系统前设备已使用时间（秒） 非必填
					"UseTimeBefore": Number($scope.form.UseTimeBefore),
					//    录入系统后设备已使用时间 非必填
					"UseTimeAfter": Number($scope.form.UseTimeAfter)
				}
			}
			//console.log(data)
			//return false;
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("修改设备",data);
            if(data.Rcode == "1000"){
				toaster.pop('success', '修改成功！');
				$modalInstance.close(true);
            }else{
              toaster.pop('warning',data.Reason);
			}
		});
	}



	//    打开弹窗  选择设备型号
	$scope.modalOpenDevice = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_device.html',
			controller: 'modalGetDeviceCtrl',
			resolve: {
				items: function() {
					return {'Type':'1','show':false};
				}
			}
		});

		modalInstance.result.then(function(deviceItem) {
			console.log("bb",deviceItem)
			if(!deviceItem){
				$scope.form.ModelId = "";
				$scope.form.ModelName = "";
			}else{
				$scope.form.ModelId = deviceItem.Id;
				$scope.form.ModelName = deviceItem.Name;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择教室
    $scope.modalOpenClassroom =function(){
		console.log("打开弹窗 -选择教室");
    	var modalInstance=$modal.open({
    		templateUrl:'../html/modal/modal_school.html',
    		controller:'modalGetClassRoomCtrl',
    		resolve:{
    			items:function(){
					return $scope.itens;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem){
				$scope.form.ClassroomId = "";
				$scope.form.ClassroomName = "";
			}else{
				if(selectedItem.addCode == "classroom"){
					$scope.form.ClassroomId = selectedItem.addId;
					$scope.form.ClassroomName = selectedItem.add;
					//   查询教室节点列表
					getNodeList();
				}else{
					$scope.form.ClassroomId = "";
					$scope.form.ClassroomName = "";
					$modal.open({
						templateUrl: 'modal/modal_alert_all.html',
						controller: 'modalAlert2Conter',
						resolve: {
              items: function () {
                return {"type":'warning',"msg":'请选择到教室！'};
							}
						}
					});
				}
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

	$scope.ok = function(){
		switch($scope.items.operate){
			case "add":
				//
				saveDevice();
				break;
			case "look":
				//
				break;
			case "edit":
				//
				saveDeviceEdit();
				break;
		}
	}

	$scope.run = function(){
		//   节点开关状态
		$scope.form.NodeSwitchStatusItems = [
			{"val":"on","title":"开"},
			{"val":"off","title":"关"}
		];
		//   接入方式
		$scope.form.JoinMethodItems = [
			{"val":"node","title":"node"},
			{"val":"pjlink","title":"pjlink"}
		];
		$scope.form.JoinMethodItem = $scope.form.JoinMethodItems[0];
		$scope.form.JoinMethod = $scope.form.JoinMethodItems[0].val;
		//
		switch($scope.items.operate){
			case "add":
				//
				break;
			case "look":
				//
				$scope.form = $.extend({},$scope.form,$scope.items.item);
				$scope.form.ClassroomName = $scope.items.item.Campusname + "-" + $scope.items.item.Buildingname + "-" + $scope.items.item.Classroomsname;
				break;
			case "edit":
				//
				$scope.form = $.extend({},$scope.form,$scope.items.item);
				$scope.form.ClassroomName = $scope.items.item.Campusname + "-" + $scope.items.item.Buildingname + "-" + $scope.items.item.Classroomsname;
				break;
		}
	}
	$scope.run();

}]);
