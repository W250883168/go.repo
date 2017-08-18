//开启路由器
app.controller("sbglJdpzKzmlContr", ['$scope', 'httpService', '$modal','toaster', function ($scope, httpService, $modal,toaster) {
    console.log("节点配置-控制命令");
    $scope.kzmlItems = [];
    $scope.page = {
        //   超始页
        "index":1,
        //   每页显示
        "oneSize":10,
        //   页码显示条数
        "pageNumber":5
    }
    //  查询条件
    $scope.form = {
        //   搜索关键词
        "KeyWord":"",
        "ModelId":""
    }
    //  查询列表
    var getNodeModelCMDList = function(pageindex,pagesize){
        Number(pageindex) > 0 ? pageindex = Number(pageindex) : pageindex = 1;
        Number(pagesize) > 0 ? pagesize = Number(pagesize) : pagesize = 10;
        var url = config.HttpUrl + "/device/getNodeModelCMDList";
        var data = {
            Auth:{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            Page: {
                PageIndex: -1
                //PageSize: pagesize
            },
            Para:{
                KeyWord:$scope.form.KeyWord,
                ModelId:$scope.form.ModelId
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("控制命令查询",data);
            if(data.Rcode =="1000"){
                $scope.kzmlItems = data.Result.Data;
                //分页
                $scope.backPage = pageFn(data.Result.Page, $scope.page.pageNumber);
                console.log('backPage',$scope.backPage);
            }else{
                console.log(data.Reason);
            }
        });
    }

    /*  -------------------- 分页、页码  -----------------------  */
    $scope.backPage = {};
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
        if(!(Number(pageindex) > 0))return false;
        if(pageindex > 0 && pageindex <= $scope.backPage.PageCount){
            getNodeModelCMDList(pageindex,$scope.page.oneSize);
        }
    }
    /*  -------------------- 分页、页码  -----------------------  */

    //  删除列表
    var deleteNodeModelCMD = function(Id){
        var url = config.HttpUrl + "/device/deleteNodeModelCMD";
        var data = {
            Auth:{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            Para:{
                Id:Id
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("删除成功",data);
            if(data.Rcode =="1000"){
                getNodeModelCMDList();
            	toaster.pop('success', '删除成功！');
            }else{
            	toaster.pop('warning', data.Reason);
            }
        });
    }
    //  添加按钮功能
    $scope.openModalAddMl = function (item,str) {
        if(!str)str = "";
        if(!item)item = {};
        var modalInstance = $modal.open({
            templateUrl: '../project/sbgl/html/sbgl/jdpz/modal_kzml.html',
            controller: 'modalSbglJdpzKzmlContr',
            windowClass: 'm-modal-sbgl-jdpz',
            resolve: {
                items: function () {
                    return {"operate":str,"item":item};
                }
            }
        });
        modalInstance.result.then(function(bul){
            console.log("bul",bul);
            if(bul){
                getNodeModelCMDList();
            }
        });
    }




    //  删除前提示
	var deleteBefore = function(Id,delFn) {
		var url = config.HttpUrl + "/device/onDeletingNodeModelCmd";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				Id: Id
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
						"msg": '确定删除当前节点控制命令吗？'
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				deleteBefore(item.Id, function(rBul) {
					//   是否有关联内容
					if(rBul){
						var modalInstance2 = $modal.open({
							templateUrl: 'modal/modal_alert_all.html',
							controller: 'modalAlert2Conter',
							resolve: {
								items: function() {
									return {
										"type": 'warning',
										"msg": '警告：此操作将影响与当前命令关联的节点控制，同时也会影响关联的设备控制！'
									};
								}
							}
						});
						modalInstance2.result.then(function(bul) {
							console.log("bul", bul);
							if(bul) {
								//    删除
								deleteNodeModelCMD(item.Id);
							}
						});
					}else{
						//  删除
						deleteNodeModelCMD(item.Id);
					}
				});
			}
		});

    }
    //   查询
    $scope.searchPost = function(){
        getNodeModelCMDList();
    }
    //   回车查询
    $scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            getNodeModelCMDList();
        }
    }
    $scope.run = function(){
        getNodeModelCMDList();
    }
    $scope.run();
}]);

//设备管理-节点配置-控制命令-添加控制命令弹窗
app.controller("modalSbglJdpzKzmlContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("添加控制命令弹窗",items);
    $scope.items = items;
    $scope.form={
        "Id":null,
        "ModelId":"",
        "NodeModelName":"",
        "CmdDescription":"",
        "CmdCode":"",
        "CmdName":"",
        "RequestURI":"",
        "URIQuery":"",
        "RequestType":""
    }


    /*************节点型号--查询列表***************/
    $scope.jdxhItems = [];
    //  查询条件
    $scope.forms = {
        //   搜索关键词
        "KeyWord":"",
        "ModelId":""
    }
    //  查询列表
    var getNodeModelList = function(){
        var url = config.HttpUrl + "/device/getNodeModelList";
        var data = {
            Auth:{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            Para:{
                KeyWord:$scope.forms.KeyWord,
                ModelId:$scope.forms.ModelId
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("节点型号查询",data);
            if(data.Rcode =="1000"){
                $scope.jdxhItems = data.Result.Data;
                for (var item in $scope.jdxhItems) {
                  if( $scope.jdxhItems[item].Id == $scope.form.ModelId){
                    $scope.form.jdxhItem = $scope.jdxhItems[item];
                  }
                }
            }else{
                console.log(data.Reason);
            }
        });
    }
    /*------节点型号--查询列表------------*/

    //控制命令--添加
    var saveNodeModelCMDAdd = function(){
	  if(!(formValidate($scope.form.CmdCode).minLength(0).outMsg(2401).isOk))return false;
      if(!(formValidate($scope.form.CmdName).minLength(0).outMsg(2403).isOk))return false;
      if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2402).isOk))return false;
      if(!(formValidate($scope.form.RequestURI).minLength(0).outMsg(2405).isOk))return false;
	  if(!(formValidate($scope.form.URIQuery).minLength(0).outMsg(2406).isOk))return false;
      if(!(formValidate($scope.form.RequestType).minLength(0).outMsg(2404).isOk))return false;
        var url = config.HttpUrl+"/device/saveNodeModelCMD";
        var data={
            Auth:{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            Para: {
                //"Id": $scope.items.Id,
                "ModelId":$scope.form.ModelId,
                "CmdCode":$scope.form.CmdCode,
                "CmdName":$scope.form.CmdName,
                "CmdDescription":$scope.form.CmdDescription,
                "RequestURI":$scope.form.RequestURI,
                "URIQuery":$scope.form.URIQuery,
                "RequestType":$scope.form.RequestType
            }
        }
        var promise = httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("控制命令-添加",data);
            if(data.Rcode == "1000"){
            	toaster.pop('success', '添加成功！');
                $modalInstance.close(true);
            }else{
            	toaster.pop('warning', data.Reason);
            }
        },function(reason){},function(update){});
    }
    //节点型号-控制命令--修改
    var saveNodeModelCMD = function(Id){
        if(!(formValidate($scope.form.CmdCode).minLength(0).outMsg(2401).isOk))return false;
        if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2402).isOk))return false;
        if(!(formValidate($scope.form.CmdName).minLength(0).outMsg(2403).isOk))return false;
        if(!(formValidate($scope.form.RequestType).minLength(0).outMsg(2404).isOk))return false;
        if(!(formValidate($scope.form.RequestURI).minLength(0).outMsg(2405).isOk))return false;
        if(!(formValidate($scope.form.URIQuery).minLength(0).outMsg(2406).isOk))return false;
        if(!Id)return false;
        var url = config.HttpUrl+"/device/saveNodeModelCMD";
        var data={
            Auth:{
                "Usersid": config.GetUser().Usersid,
                "Rolestype": config.GetUser().Rolestype,
                "Token": config.GetUser().Token,
                "Os": "WEB"
            },
            Para: {
                "Id":$scope.form.Id,
                "ModelId":$scope.form.ModelId,
                "NodeModelName":$scope.form.NodeModelName,
                "CmdDescription":$scope.form.CmdDescription,
                "CmdCode":$scope.form.CmdCode,
                "CmdName":$scope.form.CmdName,
                "RequestURI":$scope.form.RequestURI,
                "URIQuery":$scope.form.URIQuery,
                "RequestType":$scope.form.RequestType
            }
        }
        var promise=httpService.ajaxPost(url,data);
        promise.then(function(data){
            console.log("控制命令修改",data);
            if(data.Rcode == "1000"){
            	toaster.pop('success', '修改成功！');
                $modalInstance.close(true);
            }else{
            	toaster.pop('warning', data.Reason);
            }
        },function(reason){},function(update){});
    }
    //保存按钮
    $scope.ok = function(){
        if($scope.items.operate == "edit"){
            saveNodeModelCMD($scope.items.item.Id);
        }else{
            saveNodeModelCMDAdd();
        }
    }
    //	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}
  $scope.changeModelItem = function (item) {
    $scope.form.ModelId = item.Id;
  }
    $scope.run = function(){
        if(items.operate == "edit" || items.operate == "see"){
            $scope.form = $.extend({},$scope.form,$scope.items.item);
        }
        getNodeModelList();
    }
    $scope.run();
}]);
