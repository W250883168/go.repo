'use strict';
/**
 * Created by Administrator on 2016/9/7.
 */

/*   设备日志     */
app.controller('sbglSbrzContr', ['$scope', 'httpService','toaster',function($scope, httpService,toaster) {
	console.log("设备日志")

	//   设备日志list
	$scope.AllOperateLogList = [];
	//   返回分页
	$scope.backPage = {};
	//   input
	$scope.from = {
		//   开始时间
		"fromTime":"",
		//   结束时间
		"toTime":"",
		//   关键词
		"keyWord":""
	}
	//
	$scope.page = {
		//   超始页
		"index":1,
		//   每页显示
		"oneSize":15,
		//   页码显示条数
		"pageNumber":5
	}



		//   HTML ready
	$scope.showFromDate = function() {
		jeDate({
			dateCell: "#sbrz_begindate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.from.fromTime = val;
				//
				$scope.searchPost();
			},
			okfun: function(elem,val) {
				$scope.from.fromTime = val;
				//
				$scope.searchPost();
			},
			clearfun:function(elem, val) {
				$scope.from.fromTime = "";
			}
		});
	}

	$scope.showToDate = function() {
		jeDate({
			dateCell: "#sbrz_enddate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.from.toTime = val;
				//
				$scope.searchPost();
			},
			okfun: function(elem,val) {
				$scope.from.toTime = val;
				//
				$scope.searchPost();
			},
			clearfun:function(elem, val) {
				$scope.from.toTime = "";
			}
		});
	}

	//   设备日志查询

	$scope.getAllOperateLogList = function(fromtime,totime,keyword,pageindex,pagesize) {
		Number(pageindex) > 0 ? pageindex = Number(pageindex) : pageindex = 1;
		Number(pagesize) > 0 ? pagesize = Number(pagesize) : pagesize = 10;
		var url = config.HttpUrl + "/device/getAllOperateLogList";
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
				FromTime: fromtime,
				ToTime: totime,
				KeyWord: keyword
			}
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log($scope.AllOperateLogList)
			if(data.Rcode == "1000") {
				$scope.AllOperateLogList = data.Result.Data;
				//   分页
				$scope.backPage = pageFn(data.Result.Page,5);
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
	//   查询
	$scope.searchPost = function(){
		$scope.page.index = 1;
		$scope.getAllOperateLogList($scope.from.fromTime,$scope.from.toTime,$scope.from.keyWord,$scope.page.index,$scope.page.oneSize);
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
			$scope.getAllOperateLogList($scope.from.fromTime,$scope.from.toTime,$scope.from.keyWord,pageindex,$scope.page.oneSize);
		}
	}
	/*  -------------------- 分页、页码  -----------------------  */

	//   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            $scope.searchPost();
        }
	}

	//  run
	var run = function(){
		//   预置时间
		var date = new Date();
	    var seperator1 = "-";
	    var month = date.getMonth() + 1;
	    var strDate = date.getDate();
	    if (month >= 1 && month <= 9) {
	        month = "0" + month;
	    }
	    if (strDate >= 0 && strDate <= 9) {
	        strDate = "0" + strDate;
	    }

		$scope.from.fromTime = date.getFullYear() + seperator1 + month + seperator1 + strDate + " 00:00:00";
		$scope.from.toTime = date.getFullYear() + seperator1 + month + seperator1 + strDate + " 23:59:59";

		$scope.getAllOperateLogList($scope.from.fromTime,$scope.from.toTime,$scope.from.keyWord,$scope.page.index,$scope.page.oneSize);
	}
	run();

}]);
