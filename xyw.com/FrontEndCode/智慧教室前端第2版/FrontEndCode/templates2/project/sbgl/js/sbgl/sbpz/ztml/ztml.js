//开启路由器
app.controller("sbglZtmlContr", ['$scope', 'httpService', '$modal','toaster', function ($scope, httpService, $modal,toaster) {
    console.log("设备配置-状态命令");

    /**  表单关键词   */
   	$scope.form = {
   		//   关键词
   		"KeyWord":"",
   		//设备型号ID:字符串   ？？
   		"ModelId":"",
   		//    设备型号名称
   		"ModelName":""
   	}
    /**  设备状态命令列表  */
   	$scope.ModelItems = [];

    //    分页
    $scope.backPage = {
        PageIndex:1,
        PageSize:10
    }

    $scope.openModalAddMl = function (str,item) {
        var modalInstance = $modal.open({
            templateUrl: '../project/sbgl/html/sbgl/sbpz/modal_ztml.html',
            controller: 'modalZtmlContr',
            windowClass: 'm-modal-sbgl-sbpz',
            resolve: {
                items: function () {
                    return {"operate":str,"item":item};
                }
            }
        });

        modalInstance.result.then(function(bol) {
			console.log(bol)
			if(!bol){
				//
			}else{
				getDeviceModelStatusCMDList();
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
    }

    //打开窗口  设备型号
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
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}


	/**
	 * 取设备状态命令列表
	 */
    var getDeviceModelStatusCMDList = function(){
        var url = config.HttpUrl + "/device/getDeviceModelStatusCMDList";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Page": {
                "PageIndex":$scope.backPage.PageIndex,
                "PageSize":$scope.backPage.PageSize
            },
            "Para":{
                "KeyWord":$scope.form.KeyWord,
                "ModelId":$scope.form.ModelId
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("取设备状态命令列表",data);
            if(data.Rcode =="1000"){
                $scope.ModelItems = data.Result.Data;
                //分页
                $scope.backPage = pageFn(data.Result.Page,5);
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }


    /**
	 * 删除设备状态命令
	 */
    var deleteDeviceModelStatusCMD = function(ModelId,StatusCode){
        var url = config.HttpUrl + "/device/deleteDeviceModelStatusCMD";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Para":{
                "ModelId":ModelId,
                "StatusCode":StatusCode
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("删除设备状态命令",data);
            if(data.Rcode =="1000"){
            	getDeviceModelStatusCMDList();
              toaster.pop("success",'删除成功！');
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }
    
    
    
    //  删除前提示
	var deleteBefore = function(Id,Code,delFn) {
		var url = config.HttpUrl + "/device/onDeletingDeviceModelStatusCmd";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				ModelId: Id,
				StatusCode: Code
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除前提示", data);
			if(data.Rcode == "1000") {
				delFn(data.Result);
			} else {
				return false;
			}
		});
	}
    
    
    
    //  删除弹窗功能
    $scope.deleteItem = function(item){
    	
    	var modalInstance = $modal.open({
			templateUrl: 'modal/modal_alert_all.html',
			controller: 'modalAlert2Conter',
			resolve: {
				items: function() {
					return {
						"type": 'warning',
						"msg": '确定删除当前设备型号状态命令吗？'
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				deleteBefore(item.Id,item.StatusCode, function(rBul) {
					//   是否有关联内容
					if(rBul){
						var modalInstance2 = $modal.open({
							templateUrl: 'modal/modal_alert_all.html',
							controller: 'modalAlert2Conter',
							resolve: {
								items: function() {
									return {
										"type": 'warning',
										"msg": '警告：此操作将影响所有与当前命令关联的设备状态获取！'
									};
								}
							}
						});
						modalInstance2.result.then(function(bul) {
							console.log("bul", bul);
							if(bul) {
								//    删除
								deleteDeviceModelStatusCMD(item.ModelId,item.StatusCode);
							}
						});
					}else{
						//  删除
						deleteDeviceModelStatusCMD(item.ModelId,item.StatusCode);
					}
				});
			}
		}); 
    	
      
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
    $scope.pageClick = function(pageindex){
        Number(pageindex) < 1 ? pageindex = 1 : pageindex = Number(pageindex);
        $scope.backPage.PageIndex = pageindex;
        getDeviceModelStatusCMDList();
    }
    /*  -------------------- 分页、页码  -----------------------  */


    //    查询
    $scope.searchPost = function(){
		$scope.backPage.PageIndex = 1;
    	getDeviceModelStatusCMDList();
    }
    //   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            getDeviceModelStatusCMDList();
        }
	}

    $scope.run = function(){
    	getDeviceModelStatusCMDList();
    }
    $scope.run();


}]);

//设备管理-设备配置-状态命令-添加控制命令弹窗
app.controller("modalZtmlContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("设备管理-设备配置-状态命令-添加控制命令弹窗");

	//
	$scope.items = items;

	//    表单差数
	$scope.form = {
		//
		"Id":null,
		//    设备型号Id：字符串 必填
		"ModelId":"",
		//    设备型号名称
		"ModelName":"",
		//    负载命令：字符串 必填
		"Payload":"",
		//    状态命令名称：字符串 必填
		"StatusName":"",
		//    状态命令编码：字符串 必填
		"StatusCode":"",
		//    状态值匹配的字符串：字符串 必填
		"StatusValueMatchString":"",
		//    是否是开关状态：字符串 非必填
		"SwitchStatusFlag":"",
		//    开启时状态值：字符串 非必填
		"OnValue":"",
		//    关闭时状态值：字符串 非必填
		"OffValue":"",
		//    状态命令序号：数值 必填
		"SeqNo":null,
		//    从编码表取值：字符串 非必填
		"SelectValueFlag":"",
		//    状态命令是否预警：字符串 非必填
		"IsAlert":"",
		//    状态命令预警条件：字符串 非必填
		"AlertWhere":"",
		//    状态命令预警描述：字符串 非必填
		"AlertDescription":""
	}



	//打开窗口  设备型号
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

		modalInstance.result.then(function(item) {
			console.log(item)
			if(!item){
				$scope.form.ModelId = "";
				$scope.form.ModelName = "";
			}else{
				$scope.form.ModelId = item.Id;
				$scope.form.ModelName = item.Name;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}


	/**
	 * 添加设备型号状态命令
	 */
    var saveDeviceModelStatusCMD = function(){
    	//   验证
		if(!(formValidate($scope.form.StatusCode).minLength(0).outMsg(2509).isOk))return false;
		if(!(formValidate($scope.form.StatusName).minLength(0).outMsg(2510).isOk))return false;
		if(!(formValidate($scope.form.SeqNo).minLength(0).outMsg(2511).isOk))return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2500).isOk))return false;
		if(!(formValidate($scope.form.Payload).minLength(0).outMsg(2512).isOk))return false;
		if(!(formValidate($scope.form.StatusValueMatchString).minLength(0).outMsg(2513).isOk))return false;
        var url = config.HttpUrl + "/device/saveDeviceModelStatusCMD";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Para":{
                //"Id": 6,
				"ModelId": $scope.form.ModelId,
				"Payload": $scope.form.Payload,
				"StatusName": $scope.form.StatusName,
				"StatusCode": $scope.form.StatusCode,
				"StatusValueMatchString": $scope.form.StatusValueMatchString,
				"SwitchStatusFlag": $scope.form.SwitchStatusFlag,
				"OnValue": $scope.form.OnValue,
				"OffValue": $scope.form.OffValue,
				"SeqNo": Number($scope.form.SeqNo)
				//"SelectValueFlag": $scope.form.SelectValueFlag,
				//"IsAlert": $scope.form.IsAlert,
				//"AlertWhere": $scope.form.AlertWhere,
				//"AlertDescription": $scope.form.AlertDescription
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("添加设备型号状态命令",data);
            if(data.Rcode =="1000"){
              toaster.pop('success', '添加成功！');
              $modalInstance.close(true);
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }


    /**
	 * 修改设备型号状态命令
	 */
    var saveDeviceModelStatusCMDEdit = function(){
    	//   验证
		if(!(formValidate($scope.form.StatusCode).minLength(0).outMsg(2509).isOk))return false;
		if(!(formValidate($scope.form.StatusName).minLength(0).outMsg(2510).isOk))return false;
		if(!(formValidate($scope.form.SeqNo.toString()).minLength(0).outMsg(2511).isOk))return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2500).isOk))return false;
		if(!(formValidate($scope.form.Payload).minLength(0).outMsg(2512).isOk))return false;
		if(!(formValidate($scope.form.StatusValueMatchString).minLength(0).outMsg(2513).isOk))return false;
        var url = config.HttpUrl + "/device/saveDeviceModelStatusCMD";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Para":{
                "Id": Number($scope.form.Id),
				"ModelId": $scope.form.ModelId,
				"Payload": $scope.form.Payload,
				"StatusName": $scope.form.StatusName,
				"StatusCode": $scope.form.StatusCode,
				"StatusValueMatchString": $scope.form.StatusValueMatchString,
				"SwitchStatusFlag": $scope.form.SwitchStatusFlag,
				"OnValue": $scope.form.OnValue,
				"OffValue": $scope.form.OffValue,
				"SeqNo": Number($scope.form.SeqNo)
				//"SelectValueFlag": $scope.form.SelectValueFlag,
				//"IsAlert": $scope.form.IsAlert,
				//"AlertWhere": $scope.form.AlertWhere,
				//"AlertDescription": $scope.form.AlertDescription
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("修改设备型号状态命令",data);
            if(data.Rcode =="1000"){
              toaster.pop('success', '修改成功！');
              $modalInstance.close(true);
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }


	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

	//   ok
	$scope.ok = function(){
		switch($scope.items.operate){
			case "add":
				//
				saveDeviceModelStatusCMD();
			break;
			case "edit":
				//
				saveDeviceModelStatusCMDEdit();
			break;
		}

	}

	$scope.run = function(){
		switch($scope.items.operate){
			case "add":
				//
			break;
			case "look":
				//
				$scope.form = $.extend({},$scope.form,$scope.items.item);
			break;
			case "edit":
				//
				$scope.form = $.extend({},$scope.form,$scope.items.item);
			break;
		}
	}
	$scope.run();


}]);

app.controller("modalZtmlAddContr",['$scope', 'httpService', '$modal','$location','toaster',function ($scope, httpService, $modal,$location,toaster) {
	console.log("设备管理-设备配置-状态命令-添加ADD");

	//    设备型号状态编码
	$scope.statusItem = {
		//   设备型号Id
		"ModelId":$location.search().mid,
		//   状态名称
		"StatusName":$location.search().code_name,
		//   状态编码
		"StatusCode":$location.search().code
	}

	//    分页
    $scope.backPage = {
        PageIndex:1,
        PageSize:20
    }

    //    状态列表数组
    $scope.CodeList = [];


	/**
	 * 取状态值编码列表
	 */
	var getDeviceModelStatusValueCodeList = function(){
        var url = config.HttpUrl + "/device/getDeviceModelStatusValueCodeList";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Page": {
                "PageIndex":$scope.backPage.PageIndex,
                "PageSize":$scope.backPage.PageSize
            },
            "Para":{
                "KeyWord":"",
                "ModelId":$scope.statusItem.ModelId.toString(),
                "StatusCode":$scope.statusItem.StatusCode
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("取状态值编码列表",data);
            if(data.Rcode =="1000"){
                $scope.CodeList = data.Result.Data;
                //分页
                $scope.backPage = pageFn(data.Result.Page,5);
            }else{
                toaster.pop('warning',data.Reason);
            }
        });
    }


	/**
	 * 状态值编码删除
	 */
	var deleteDeviceModelStatusValueCode = function(Id){
		if(!Id)return false;
        var url = config.HttpUrl + "/device/deleteDeviceModelStatusValueCode";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Para":{
                "Id":Number(Id)
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            if(data.Rcode =="1000"){
            	getDeviceModelStatusValueCodeList();
              toaster.pop('success', '删除成功！');
            }else{
              toaster.pop('warning',data.Reason);
            }
        });
    }


	/**
	 * 删除状态值编码
	 */
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'你确定要删除吗？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        deleteDeviceModelStatusValueCode(item.Id);
      }
    });
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
			getDeviceModelStatusValueCodeList();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */


	//打开窗口  设备型号
	$scope.modalOpenCode = function() {
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbpz/ztml/code/modal_code.html',
			controller: 'modalSbglSbpzZtmlCodeCtrl',
			resolve: {
				items: function() {
					return $scope.statusItem;
				}
			}
		});

		modalInstance.result.then(function(bol) {
			console.log(bol)
			if(!bol){
				//
			}else{
				getDeviceModelStatusValueCodeList();
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}



//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}



	$scope.run = function(){
		getDeviceModelStatusValueCodeList();
		//   面包屑标题
		$scope.operatetitle = "状态值设置";
	}
	$scope.run();

}]);



/**
 * 添加状态值编码弹窗
 */
app.controller("modalSbglSbpzZtmlCodeCtrl",['$scope','httpService','$modalInstance','items','formValidate','toaster',function ($scope, httpService,$modalInstance,items,formValidate,toaster) {
	console.log("设备管理-设备配置-状态命令-添加ADD");

	$scope.items = items;

	//    设备型号状态编码
	$scope.form = {
		//   id
		"Id":null,
		//   设备型号Id：字符串 必填
		"ModelId":"",
		//   状态命令编码：字符串 必填
		"StatusCode":"",
		//   状态编码：字符串 必填
		"StatusValueCode":"",
		//   状态名称：字符串 必填
		"StatusValueName":"",
		//   是否预警：字符串 非必填
		"IsAlert":""
	}

	/**
	 * 添加状态值编码
	 */
	var saveDeviceModelStatusValueCode = function(){
		if(!(formValidate($scope.form.StatusValueCode).minLength(0).outMsg(2514).isOk))return false;
		if(!(formValidate($scope.form.StatusValueName).minLength(0).outMsg(2515).isOk))return false;
        var url = config.HttpUrl + "/device/saveDeviceModelStatusValueCode";
        var data = {
            "Auth":{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            "Para":{
                //"Id":Number($scope.form.Id),
                "ModelId":$scope.items.ModelId,
                "StatusCode":$scope.items.StatusCode,
                "StatusValueCode":$scope.form.StatusValueCode,
                "StatusValueName":$scope.form.StatusValueName
                //"IsAlert":$scope.form.IsAlert
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("添加状态值编码",data);
            if(data.Rcode =="1000"){
              toaster.pop('success', '添加成功！');
              $modalInstance.close(true);
            }else{
              toaster.pop('warning', data.Reason);
            }
        });
    }

	//  ok
	$scope.ok = function(){
		saveDeviceModelStatusValueCode();
	}


	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

}]);
