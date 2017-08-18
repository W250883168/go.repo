//开启路由器
app.controller("sbglSbpzKzmlContr", ['$scope', 'httpService', '$modal','toaster', function ($scope, httpService, $modal,toaster) {
    console.log("设备配置-控制命令");
    $scope.SbpzKzmlItems = [];
	$scope.page = {
		//   超始页
		"index":1,
		//   每页显示
		"oneSize":10,
		//   页码显示条数
		"pageNumber":5
	}
	$scope.form = {
		"KeyWord":"",
		"ModelId":"",
		"ModelName":""
	}
	//	查询列表
	var getDeviceModelControlCMDList = function(pageindex,pagesize){
		Number(pageindex) > 0 ? pageindex = Number(pageindex) : pageindex = 1;
		Number(pagesize) > 0 ? pagesize = Number(pagesize) : pagesize = 10;
		var url = config.HttpUrl + "/device/getDeviceModelControlCMDList";
		var data = {
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: pageindex,
				PageSize: pagesize
			},
			Para:{
				KeyWord:$scope.form.KeyWord,
				ModelId:$scope.form.ModelId
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("控制命令查询列表",data);
			if(data.Rcode == "1000"){
				$scope.SbpzKzmlItems = data.Result.Data;
				//分页
				$scope.backPage = pageFn(data.Result.Page, $scope.page.pageNumber);
			}else{
        toaster.pop('warning',data.Reason);
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
			getDeviceModelControlCMDList(pageindex,$scope.page.oneSize);
		}
	}
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
	var deleteDeviceModelControlCMD = function(Id){
		var url = config.HttpUrl + "/device/deleteDeviceModelControlCMD";
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
				getDeviceModelControlCMDList();
        toaster.pop('success', '删除成功！');
			}else{
        toaster.pop('warning',data.Reason);
			}
		});
	}
  //  删除弹窗功能
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'确定删除此设备型号控制命令？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        deleteDeviceModelControlCMD(item.Id);
      }
    });
  }
	//  添加按钮功能
	$scope.openModalAddMl = function (item,str) {
		if(!str)str = "";
		if(!item)item = {};
		var modalInstance = $modal.open({
			templateUrl: '../project/sbgl/html/sbgl/sbpz/modal_kzml.html',
			controller: 'modalSbglSbpzKzmlContr',
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

//设备管理-设备配置-控制命令-添加控制命令弹窗
app.controller("modalSbglSbpzKzmlContr",['$scope', 'httpService', '$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService, $modal,$modalInstance,items,formValidate,toaster) {
	console.log("设备管理-设备配置-控制命令-添加控制命令弹窗");
	$scope.items = items;
	console.log("弹窗",$scope.items);
	$scope.form = {
		"Id":null,
		"ModelIdName":"",
		"ModelId":"",
		"ModelName":"",
		"CmdCode":"",
		"CmdName":"",
		"RequestURI":"",
		"URIQuery":"",
		"CmdDescription":"",
		"RequestType":"",
		"Payload":"",
		"DelayMillisecond":null,

		"CloseCmdFlag":"",
		//  ----  整体操作打开命令   选中
		"CloseCmdFlagItem":"",
		//  ----  整体操作打开命令   是(1)  否(0)
		"CloseCmdFlagItems":[],

		"OpenCmdFlag":"",
		//  ----  整体操作关闭命令   选中
		"OpenCmdFlagItem":"",
		//  ----  整体操作关闭命令   是(1)  否(0)
		"OpenCmdFlagItems":[]
	}
	//  ********查询列表*******************
	$scope.jdxhItems = [];
	//  查询条件
	$scope.forms = {
		"KeyWord":"",
		"ModelId":""
	}


	//    整体操作打开命令
	$scope.changeCloseCmdFlagItem = function(item){
		$scope.form.CloseCmdFlagItem = item;
		$scope.form.CloseCmdFlag = item.val;
	}

	//    整体操作关闭命令
	$scope.changeOpenCmdFlagItem = function(item){
		$scope.form.OpenCmdFlagItem = item;
		$scope.form.OpenCmdFlag = item.val;
	}


	var getDeviceModelList = function(){
		var url = config.HttpUrl + "/device/getDeviceModelList";
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
			console.log("设备型号查询",data);
			if(data.Rcode =="1000"){
				$scope.jdxhItems = data.Result.Data;
			}else{
        toaster.pop('warning',data.Reason);
			}
		});
	}
	/*------节点型号--查询列表------------*/
	//控制命令--添加
	var saveDeviceModelControlCMDAdd = function(){
		if(!(formValidate($scope.form.CmdCode).minLength(0).outMsg(2503).isOk))return false;
		if(!(formValidate($scope.form.CmdName).minLength(0).outMsg(2504).isOk))return false;
    if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2500).isOk))return false;
    if(!(formValidate($scope.form.RequestURI).minLength(0).outMsg(2505).isOk))return false;
    if(!(formValidate($scope.form.URIQuery).minLength(0).outMsg(2506).isOk))return false;
    if(!(formValidate($scope.form.RequestType).minLength(0).outMsg(2508).isOk))return false;
    if(!(formValidate($scope.form.Payload).minLength(0).outMsg(2507).isOk))return false;

		var url = config.HttpUrl+"/device/saveDeviceModelControlCMD";
		var data={
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				"Id": 0,
				"ModelId":$scope.form.ModelId,
				"CmdCode":$scope.form.CmdCode,
				"CmdName":$scope.form.CmdName,
				"RequestURI":$scope.form.RequestURI,
				"URIQuery":$scope.form.URIQuery,
				"CmdDescription":$scope.form.CmdDescription,
				"RequestType":$scope.form.RequestType,
				"Payload":$scope.form.Payload,
				"DelayMillisecond":Number($scope.form.DelayMillisecond),
				"CloseCmdFlag":$scope.form.CloseCmdFlag,
				"OpenCmdFlag":$scope.form.OpenCmdFlag
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
	//控制命令--修改
	var saveDeviceModelControlCMD = function(Id){
		if(!(formValidate($scope.form.CmdCode).minLength(0).outMsg(2503).isOk))return false;
		if(!(formValidate($scope.form.CmdName).minLength(0).outMsg(2504).isOk))return false;
		if(!(formValidate($scope.form.RequestURI).minLength(0).outMsg(2505).isOk))return false;
		if(!(formValidate($scope.form.URIQuery).minLength(0).outMsg(2506).isOk))return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2500).isOk))return false;
		if(!(formValidate($scope.form.Payload).minLength(0).outMsg(2507).isOk))return false;
		if(!(formValidate($scope.form.RequestType).minLength(0).outMsg(2508).isOk))return false;
		if(!Id)return false;
		var url = config.HttpUrl+"/device/saveDeviceModelControlCMD";
		var data={
			Auth:{
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				"Id": Number($scope.form.Id),
				"ModelId":$scope.form.ModelId,
				"ModelIdName":$scope.form.ModelIdName,
				"CmdCode":$scope.form.CmdCode,
				"CmdName":$scope.form.CmdName,
				"RequestURI":$scope.form.RequestURI,
				"URIQuery":$scope.form.URIQuery,
				"CmdDescription":$scope.form.CmdDescription,
				"RequestType":$scope.form.RequestType,
				"Payload":$scope.form.Payload,
				"DelayMillisecond":Number($scope.form.DelayMillisecond),
				"CloseCmdFlag":$scope.form.CloseCmdFlag,
				"OpenCmdFlag":$scope.form.OpenCmdFlag
			}
		}
		var promise = httpService.ajaxPost(url,data);
		promise.then(function(data){
			console.log("控制命令-修改",data);
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
			saveDeviceModelControlCMD($scope.items.item.Id);
		}else{
			saveDeviceModelControlCMDAdd();
		}
	}
	//	取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
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

	$scope.run = function(){
		//   整体操作打开命令
		$scope.form.CloseCmdFlagItems = [
			{"val":1,"title":"是(1)"},
			{"val":0,"title":"否(0)"}
		];

		//   整体操作关闭命令
		$scope.form.OpenCmdFlagItems = $scope.form.CloseCmdFlagItems;


		if(items.operate == "edit" || items.operate == "see"){
			$scope.form = $.extend({},$scope.form,$scope.items.item);
			$scope.form.ModelName =$scope.items.item.ModelIdName;
			//  整体操作打开命令
			if($scope.form.CloseCmdFlag == '0'){
				$scope.form.CloseCmdFlagItem = $scope.form.CloseCmdFlagItems[1];
			}else{
				$scope.form.CloseCmdFlagItem = $scope.form.CloseCmdFlagItems[0];
			}
			//  整体操作关闭命令
			if($scope.form.OpenCmdFlag == '0'){
				$scope.form.OpenCmdFlagItem = $scope.form.OpenCmdFlagItems[1];
			}else{
				$scope.form.OpenCmdFlagItem = $scope.form.OpenCmdFlagItems[0];
			}
		}
		getDeviceModelList();
	}
	$scope.run();
}]);
