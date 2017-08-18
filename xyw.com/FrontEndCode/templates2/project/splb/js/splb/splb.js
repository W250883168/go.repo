'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */
app.controller("splbindexContr",['$scope','$state','httpService','$modal','$timeout',function($scope,$state,httpService,$modal,$timeout){
	console.log("视频录播");
	//   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:10
    }
	//  form
	$scope.form = {
		"KeyWords":""
	}
	//   
	$scope.videoList = [];
	/**
	 * 取视频列表
	 */
	var getvideolist = function() {
		var url = config.HttpUrl + "/vod/getvideolist";
		var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "KeyWords":$scope.form.KeyWords,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取视频列表",data)
			if(data.Rcode == "1000") {
				$scope.videoList = data.Result.PageData;
				//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);	
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;	
				}
				$scope.backPage = pageFn(objPage,5);
            }else if(data.Rcode=="1002"){
            	$scope.videoList = [];
            	//   分页
				var objPage={PageCount:0,PageIndex:1,PageSize:10,RecordCount:0};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);	
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;	
				}
				$scope.backPage = pageFn(objPage,5);
            }else{
            	console.log(data.Reason);
            }
		}, function(reason) {}, function(update) {});
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
			getvideolist();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */
    
    //   搜索
    $scope.searchPost = function(){
		$scope.backPage.PageIndex = 1;
    	getvideolist();
    }
    
    //   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            getvideolist();
        }
	}
	
	
	
	$scope.run = function(){
		getvideolist();
	}
	$scope.run();
	
}]);


