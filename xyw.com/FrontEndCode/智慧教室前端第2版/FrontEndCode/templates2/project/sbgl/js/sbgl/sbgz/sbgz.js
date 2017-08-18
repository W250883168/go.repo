'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */

/*   设备故障     */
app.controller('sbglSbgzContr', ['$scope', 'httpService', '$modal', '$interval','$timeout','toaster', function($scope, httpService, $modal, $interval,$timeout,toaster) {
	console.log("设备故障")
	
	$scope.allAlertInfoList = [];
	//   
	$scope.page = {
		//   超始页
		"index":1,
		//   每页显示
		"oneSize":15,
		//   页码显示条数
		"pageNumber":5
	}
	$scope.searchData = {
		//  位置类型
		"SiteType":"",
		//  设备位置id：字符串
		"SiteId":"",
		//   设备型号id：字符串
		"ModelId":"",
		//   关键词
		"KeyWord":""
	}
	//   返回 位置数据
	$scope.backAdd = {};
	//   位置HTML显示
	$scope.addAdd = "";
	//   设备HTML显示
	$scope.deviceText = "";
	
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
			console.log(selectedItem)
			if(selectedItem.addId == ""){
				$scope.addAdd = "";
				$scope.backAdd = selectedItem;
				$scope.searchData.SiteType = "";
				$scope.searchData.SiteId = "";
				$scope.addAdd = "";
			}else{
				$scope.backAdd = selectedItem;
				$scope.searchData.SiteType = $scope.backAdd.addCode;
				$scope.searchData.SiteId = $scope.backAdd.addId;
				$scope.addAdd = $scope.backAdd.add;
			}
			//
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
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
				$scope.deviceText = "";
				$scope.searchData.ModelId = "";
			}else{
				$scope.searchData.ModelId = deviceItem.Id;
				$scope.deviceText = deviceItem.Name;
			}
			//
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	

	//   获取设备故障信息
	$scope.getAllAlertInfoList = function(sitetype,siteid, modelid, keyword, pageindex, pagesize) {
		Number(pageindex) > 0 ? pageindex = Number(pageindex) : pageindex = 1;
		Number(pagesize) > 0 ? pagesize = Number(pagesize) : pagesize = 10;
		var url = config.HttpUrl + "/device/getAllFaultInfoList";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Os": "WEB",
				"Token": config.GetUser().Token
			},
			Page: {
				PageIndex: pageindex,
				PageSize: pagesize
			},
			Para: {
				SiteType: sitetype,
				SiteId: siteid.toString(),
				ModelId: modelid,
				KeyWord: keyword
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log(data)
			if(data.Rcode == "1000") {
				$scope.allAlertInfoList = data.Result.Data;
				//   分页
				$scope.backPage = pageFn(data.Result.Page, $scope.page.pageNumber);
			} else {
				console.log(data.Reason);
			}
			//console.log("设备预警查询");
			//console.log(data);
		}, function(reason) {}, function(update) {});
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
			$scope.getAllAlertInfoList($scope.searchData.SiteType,$scope.searchData.SiteId,$scope.searchData.ModelId,$scope.searchData.KeyWord,pageindex,$scope.page.oneSize);
		}
	}
	/*  -------------------- 分页、页码  -----------------------  */
	
	//   查询
	$scope.searchPost = function(){
		$scope.page.index=1;
		$scope.getAllAlertInfoList($scope.searchData.SiteType,$scope.searchData.SiteId,$scope.searchData.ModelId,$scope.searchData.KeyWord,$scope.page.index,$scope.page.oneSize);
	}
	
	//   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            $scope.searchPost();
        }
	}
	
	
	
	/* -------- 添加  -------- */
	
	
	
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
				toaster.pop('success',"删除成功！");
				$scope.getAllAlertInfoList($scope.searchData.SiteType,$scope.searchData.SiteId,$scope.searchData.ModelId,$scope.searchData.KeyWord,$scope.page.index,$scope.page.oneSize);
			} else {
				toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	
	
	//   故障删除
	$scope.delFault = function(item,str){
		if(str == 'delete'){
			var modalInstance = $modal.open({
		      templateUrl: 'modal/modal_alert_all.html',
		      controller: 'modalAlert2Conter',
		      resolve: {
		        items: function () {
		          return {"type":'warning',"msg":'确定删除？'};
		        }
		      }
		    });
		    modalInstance.result.then(function(bul){
		      if(bul){
		        deleteFault(item.FaultId);
		      }
		    });
		}
	}
	
	
	
	//   故障添加打开弹窗
	$scope.modalOpenAddFault = function(item,str) {
		console.log(item)
		
		//   加入ID
		if(!item){
			
		}else{
			item.Id = item.FaultId;
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
				$scope.getAllAlertInfoList($scope.searchData.SiteType,$scope.searchData.SiteId,$scope.searchData.ModelId,$scope.searchData.KeyWord,$scope.page.index,$scope.page.oneSize);
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
	
	
	
	
	
	/* -------- 添加 End -------- */
	
	
	
	//   run
	var run = function(){
		$scope.getAllAlertInfoList($scope.searchData.SiteType,$scope.searchData.SiteId,$scope.searchData.ModelId,$scope.searchData.KeyWord,$scope.page.index,$scope.page.oneSize);
	}
	run();

}]);


