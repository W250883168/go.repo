'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 课程管理-课程计划-章节配置列表
 */


/*    课程管理-课程计划-章节配置列表      */
app.controller("kcglKcjhZjpzContr", ['$scope', 'httpService', '$modal','$location','toaster', function ($scope, httpService, $modal,$location,toaster) {
    console.log("课程管理-课程计划-章节配置列表")

    //     课程ID
    $scope.Curriculumsid = $location.search().ccid;
    //     课程名称
    $scope.Curriculumname = $location.search().ccname;
    //     课程ID
    $scope.CurriculumsclasscentreId = $location.search().cctid;

    $scope.form = {
    	//   中间ID
    	"Curriculumsclasscentreid":null,
    	//   班级ID
    	"Chaptersid":null,
    	//   班级名称
    	"Chaptersname":"",
    	//   上课老师ID
    	"Usersid":null,
    	//   上课老师名称
    	"Usersname":"",
    	//   上课开始时间
    	"Begindate":"",
    	//   上课结束时间
    	"Enddate":"",
    	//   教室ID
    	"Classroomid":null,
    	//   考察名称
    	"Classroomname":"",
    	//   是否直播
    	"Islive":0,
    	//   是否点播
    	"Isondemand":0
    }

    //   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:15
    }

    //   章节列表
    $scope.zjpzItems = [];


    //    打开章节配置弹窗
    $scope.openModalZjpz = function (str,item) {
    	if(!str)str = "";
    	if(!item)item = null;
        var modalInstance = $modal.open({
            templateUrl: '../project/kcgl/html/kcgl/kcjh/modal_zjpz.html',
            controller: 'modalZjpzContr',
            windowClass: 'm-modal-kcgl-zjpz',
            resolve: {
                items: function () {
                    return {"operation":str,"item":item};
                }
            }
        });

        modalInstance.result.then(function(bul) {
			console.log(bul)
			if(bul){
				curriculumclassroomchaptercentrelist();
			}

		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});

    }


    //   章节配置列表
	var curriculumclassroomchaptercentrelist = function(){
		var url=config.HttpUrl+"/system/us/curriculumclassroomchaptercentrelist";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"CurriculumsclasscentreId": Number($scope.CurriculumsclasscentreId),
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("章节配置列表",data)
			if(data.Rcode == "1000") {
				$scope.zjpzItems = data.Result.PageData;

				//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);

			}else if(data.Rcode=="1002"){
            	$scope.zjpzItems = [];
            	//   分页
				var objPage={PageCount:0,PageIndex:1,PageSize:10,RecordCount:0};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
            }else{
              toaster.pop("warning",data.Reason);
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
			curriculumclassroomchaptercentrelist();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */



    //   章节配置删除
	var curriculumclassroomchaptercentredel = function(Id){
		if(!Id)return false;
		var url=config.HttpUrl+"/system/us/curriculumclassroomchaptercentredel";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Number(Id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("章节配置列表",data)
			if(data.Rcode == "1000") {
				curriculumclassroomchaptercentrelist();
        toaster.pop('success',"删除成功！");
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

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
        curriculumclassroomchaptercentredel(item.CurriculumclassroomchaptercentreId);
      }
    });
  }

    //    查询
    $scope.searchPost = function(){
		$scope.backPage.PageIndex=1;
    	curriculumclassroomchaptercentrelist();
    }
    //   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            curriculumclassroomchaptercentrelist();
        }
	}




	$scope.run = function(){
		curriculumclassroomchaptercentrelist();
		//   面包屑标题
		$scope.operatetitle = "章节配置";
	}
	$scope.run();


}]);


/*    课程管理-课程计划-章节配置列表-章节配置弹窗      */
app.controller("modalZjpzContr", ['$scope', 'httpService','$modal','$modalInstance','$location','items','formValidate','toaster',function ($scope, httpService,$modal,$modalInstance,$location,items,formValidate,toaster) {
    console.log("课程管理-课程计划-章节配置列表-章节配置弹窗");

    //   章节
    $scope.items = items;

    //     课程ID
    $scope.Curriculumsid = $location.search().ccid;
    //     课程名称
    $scope.Curriculumname = $location.search().ccname;
    //     课程班级中间ID
    $scope.CurriculumsclasscentreId = $location.search().cctid;

    $scope.form = {
    	//   课程班级中间ID
    	"Curriculumsclasscentreid":$scope.CurriculumsclasscentreId,
    	//   章节ID
    	"Chaptersid":null,
    	//   章节
    	"ChaptersItem":null,
    	//   章节列表
    	"chapterslist":[],
    	//   开始时间
    	"Begindate":"",
    	//   开始时间Html 显示日期
    	"BegindateHtml":"",
    	//   时间段
    	"dateItems":[],
    	//   选中时间段
    	"dateItem":"",
    	//   结束时间
    	"Enddate":"",
    	//   上课教室
    	"Classroomid":null,
    	//   上课教室
    	"Classroomname":"",
    	//   上课老师
    	"TeacherID":null,
    	//   上课老师
    	"Nickname":"",
    	//   是否直播
    	"Islive":1,
    	// ------ 是否直播[0:否,1:是]
    	"IsliveItem":{"val":0,"title":"否"},
    	// ------ 是否直播[0:否,1:是]
    	"IsliveItems":[
    		{"val":0,"title":"否"},
    		{"val":1,"title":"是"}
    	],
    	//   是否录播
    	"Isondomian":1,
    	// ----- 是否录播[0:否,1:是]
    	"IsondemandItem":{"val":0,"title":"否"},
    	// ----- 是否录播[0:否,1:是]
    	"IsondemandItems":[
    		{"val":0,"title":"否"},
    		{"val":1,"title":"是"}
    	],
    	//   章节班级课程中间ID == CurriculumclassroomchaptercentreId
    	"Id":null
    }

    //    放入课程时间， 为数组 ， 用于select
    angular.forEach(config.dateTable,function(v,n){
    	$scope.form.dateItems.push({"val":v,"title":v,"name":n});
    });
    $scope.form.dateItem = $scope.form.dateItems[0];

    //   是否录播
    $scope.changeIsondemandItem = function(item){
    	$scope.form.IsondemandItem = item;
    	$scope.form.Isondemand = item.val;
    }
    //   是否直播
    $scope.changeIsliveItem = function(item){
    	$scope.form.IsliveItem = item;
    	$scope.form.Islive = item.val;
    }


    //    选择章节
	$scope.selectChapters = function(){
		console.log($scope.form.ChaptersItem)
		$scope.form.Chaptersid = $scope.form.ChaptersItem.Id;
	}

	//   选择时间段
	$scope.selectDataItem = function(item){
		$scope.form.dateItem = item;
		$scope.form.Begindate = $scope.form.BegindateHtml + " " + item.val.substr(0,5) + ":00";
		$scope.form.Enddate = $scope.form.BegindateHtml + " " + item.val.substr(6,5) + ":00";
	}


    //    打开弹出 -选择日期
    $scope.showFromDate = function() {
		jeDate({
			dateCell: "#kcpz_begindate",
			format: "YYYY-MM-DD",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.BegindateHtml = val;
				if($scope.form.dateItem){
					$scope.form.Begindate = $scope.form.BegindateHtml + " " + $scope.form.dateItem.val.substr(0,5) + ":00";
					$scope.form.Enddate = $scope.form.BegindateHtml + " " + $scope.form.dateItem.val.substr(6,5) + ":00";
				}
			},
			okfun: function(elem,val) {
				$scope.form.BegindateHtml = val;
				if($scope.form.dateItem){
					$scope.form.Begindate = $scope.form.BegindateHtml + " " + $scope.form.dateItem.val.substr(0,5) + ":00";
					$scope.form.Enddate = $scope.form.BegindateHtml + " " + $scope.form.dateItem.val.substr(6,5) + ":00";
				}
			},
			clearfun:function(elem, val) {
				$scope.form.BegindateHtml = "";
				//   清空时间与选择时间段下拉
				$scope.form.Begindate = "";
				$scope.form.Enddate = "";
				$scope.form.dateItem = "";
			}
		});
	}


    $scope.showToDate = function() {
		jeDate({
			dateCell: "#kcpz_enddate",
			format: "YYYY-MM-DD",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.Enddate = val;
			},
			okfun: function(elem,val) {
				$scope.form.Enddate = val;
			},
			clearfun:function(elem, val) {
				$scope.form.Enddate = "";
			}
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
				$scope.form.Classroomid = null;
				$scope.form.Classroomname = "";
			}else{
				if(selectedItem.addCode == "classroom"){
					$scope.form.Classroomid = selectedItem.addId;
					$scope.form.Classroomname = selectedItem.add;
				}else{
					$scope.form.Classroomid = null;
					$scope.form.Classroomname = "";
				}
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
  	 };

    //    打开弹窗  选择老师
		$scope.modalOpenTeacher =function(){
    	console.log("打开弹窗 -选择老师");
    	var modalInstance=$modal.open({
    		templateUrl:'../html/modal/modal_teacher.html',
    		controller:'modalGetTeacherCtrl',
    		resolve:{
    			items:function(){
    				return $scope.itens;
    			}
    		}
    	});

    	modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem){
				$scope.form.TeacherID = null;
				$scope.form.Nickname = "";
				$scope.form.Truename = "";
			}else{
				$scope.form.TeacherID = selectedItem.Usersid;
				$scope.form.Nickname = selectedItem.Nickname;
				$scope.form.Truename = selectedItem.Truename;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	};

	//    查找章节
	var chapterslist = function(Curriculumsid){
		if(!Curriculumsid) return false;
		var url = config.HttpUrl+"/system/us/chapterslist";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Curriculumsid": Number(Curriculumsid),
			"PageIndex": -1
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("查找章节",data)
			if(data.Rcode == "1000") {
				$scope.form.chapterslist = data.Result.PageData;
			} else {
        toaster.pop("warning",data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//     添加课程章节配置
	var curriculumclassroomchaptercentreadd = function(){
		if(!(formValidate($scope.form.Chaptersid).isNumber().outMsg(2904).isOk))return false;

		var url=config.HttpUrl+"/system/us/curriculumclassroomchaptercentreadd";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Curriculumsclasscentreid":Number($scope.form.Curriculumsclasscentreid),
			"Chaptersid": Number($scope.form.Chaptersid),
			"TeacherID": Number($scope.form.TeacherID),
			"Classroomid": Number($scope.form.Classroomid),
			"Begindate": $scope.form.Begindate,
			"Enddate": $scope.form.Enddate,
			"Islive": Number($scope.form.Islive),
			"Isondomian": Number($scope.form.Isondomian)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("添加课程章节配置",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '添加成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
	}

	//     修改课程章节配置
	var curriculumclassroomchaptercentrechange = function(){
		if(!(formValidate($scope.form.Chaptersid).isNumber().outMsg(2904).isOk))return false;

		var url=config.HttpUrl+"/system/us/curriculumclassroomchaptercentrechange";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            //"Curriculumsclasscentreid":Number($scope.form.Curriculumsclasscentreid),
			"Chaptersid": Number($scope.form.Chaptersid),
			"TeacherID": Number($scope.form.TeacherID),
			"Classroomid": Number($scope.form.Classroomid),
			"Begindate": $scope.form.Begindate,
			"Enddate": $scope.form.Enddate,
			"Islive": Number($scope.form.Islive),
			"Isondomian": Number($scope.form.Isondomian),
			"Id":Number($scope.form.Id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("修改课程章节配置",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '修改成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
	}


	$scope.ok = function(){
		if($scope.items.operation == "edit"){
			curriculumclassroomchaptercentrechange();
		}else{
			curriculumclassroomchaptercentreadd();
		}
	}

	//取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}


	$scope.run = function(){
		chapterslist($scope.Curriculumsid);

		switch($scope.items.operation){
			case "add":
				//
			break;
			case "edit":
				$scope.form.ChaptersItem = $scope.items.item;
				$scope.form.Chaptersid = $scope.items.item.Chaptersid;
				$scope.form.Begindate = $scope.items.item.Begindate;
				$scope.form.BegindateHtml = $scope.items.item.Begindate.substr(0,10);
				$scope.form.Enddate = $scope.items.item.Enddate;
				$scope.form.Classroomid = $scope.items.item.Classroomid;
				$scope.form.Classroomname = $scope.items.item.Classroomsname;
				$scope.form.TeacherID = $scope.items.item.TeacherId;
				$scope.form.Nickname = $scope.items.item.Truename;
				$scope.form.Truename = $scope.items.item.Truename;
				$scope.form.Islive = $scope.items.item.Islive;
				$scope.form.Isondomian = $scope.items.item.Isondomian;
				$scope.form.Id = $scope.items.item.CurriculumclassroomchaptercentreId;

				//   载入时间段
				if($scope.form.Begindate && $scope.form.Enddate){
					var dataStr = $scope.form.Begindate.substr(11,5) + "-" + $scope.form.Enddate.substr(11,5);
					for(var a in $scope.form.dateItems){
						if($scope.form.dateItems[a].val == dataStr){
							$scope.form.dateItem = $scope.form.dateItems[a];
						}
					}
				}

			break;
		}
	}
	$scope.run();

}]);

