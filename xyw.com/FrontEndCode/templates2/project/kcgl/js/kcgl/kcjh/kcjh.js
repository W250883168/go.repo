'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 课程管理-查课表
 */

app.filter('propsFilter', function() {
	return function(items, props) {
		var out = [];
		if(angular.isArray(items)) {
			items.forEach(function(item) {
				var itemMatches = false;
				var keys = Object.keys(props);
				for(var i = 0; i < keys.length; i++) {
					var prop = keys[i];
					try {
						if(typeof(props[prop]) == "number") {
							var text = props[prop];
							if(item[prop].toString() == text.toString()) {
								itemMatches = true;
								break;
							}
						} else {
							var text = props[prop].toString().toLowerCase();
							if(item[prop].toString().toLowerCase().indexOf(text) !== -1) {
								itemMatches = true;
								break;
							}
						}
					} catch(e) {}
				}
				if(itemMatches) {
					out.push(item);
				}
			});
		} else {
			out = items;
		}
		return out;
	};
});

/*    课程管理-课程计划      */
app.controller("kcglKcjhContr", ['$scope', 'httpService', '$modal','toaster', function ($scope, httpService, $modal,toaster) {
    console.log("课程管理-课程计划")

    $scope.kcjhItems = [];
    //   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:10
    }
    //  查询条件
    $scope.form = {
    	//   搜索关键词
    	"Seacrchtxt":"",
    	//   课程ID
    	"Curriculumsid":null,
    	//   课程名称
    	"CurriculumsName":"",
    	//   教师ID
    	"TeacherId":null,
    	//   教师名称
    	"TeacherName":"",
    	//   班级ID
    	"Classesid":null,
    	//   学科编码
    	"Subjectcode":""
    }

    //   课程班级配置查询
	var curriculumsclasscentrelist = function(){
		var url=config.HttpUrl+"/system/us/curriculumsclasscentrelist";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "TeacherId":Number($scope.form.TeacherId),
            "Curriculumsid":Number($scope.form.Curriculumsid),
			"Seacrchtxt": $scope.form.Seacrchtxt,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("课程班级配置查询",data)
			if(data.Rcode == "1000") {
				$scope.kcjhItems = data.Result.PageData;

				//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
        toaster.pop("success","查询成功！");
			} else {
				$scope.kcjhItems = [];
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//    课程删除
	var curriculumsclasscentredel = function(Id){
		var url=config.HttpUrl+"/system/us/curriculumsclasscentredel";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Number(Id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("课程删除",data)
			if(data.Rcode == "1000") {
				curriculumsclasscentrelist();
        toaster.pop("success","删除成功！");
			} else {
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
			curriculumsclasscentrelist();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */

    //    查询
    $scope.searchPost = function(){
		$scope.backPage.PageIndex=1;
    	curriculumsclasscentrelist();
    }
    //   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            curriculumsclasscentrelist();
        }
	}

    //    打开课程配置弹窗
    $scope.openModalKcpz = function (item,str) {
    	if(!str)str = "";
    	if(!item)item = {};
        var modalInstance = $modal.open({
            templateUrl: '../project/kcgl/html/kcgl/kcjh/modal_kcpz.html',
            controller: 'modalKcpzContr',
            windowClass: 'm-modal-kcgl-kcpz',
            resolve: {
                items: function () {
                    return {"operate":str,"item":item};
                }
            }
        });

        modalInstance.result.then(function(bul) {
			console.log(bul)
			if(bul){
				curriculumsclasscentrelist();
			}

		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});

    }


    //  打开窗口-选择老师
	$scope.modalOpenTeacher=function(){
		console.log("打开窗口-选择老师");
		var modalInstance=$modal.open({
			templateUrl:'../html/modal/modal_teacher.html',
			controller:'modalGetTeacherCtrl',
			resolve:{
				items:function(){
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if('Usersid' in selectedItem){
				$scope.form.TeacherName = selectedItem.Truename;
				$scope.form.TeacherId = selectedItem.Usersid;
			}else{
				$scope.form.TeacherName = "";
				$scope.form.TeacherId = null;
			}
			//
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}

	//    打开弹窗  选择 课程选择 - 学科筛选课程
	$scope.modalOpenClassroomCu = function() {
		console.log("打开弹窗  选择 课程选择 - 学科筛选课程");
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_course_science.html',
			windowClass: 'modal-kcgl-kcjh-add',
			controller: 'modalGetCourseCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				},
				deps: ['$ocLazyLoad',
					function($ocLazyLoad) {
						return $ocLazyLoad.load(['ui.select']).then();
					}
				]
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if('CurriculumsId' in selectedItem){
				$scope.form.CurriculumsName = selectedItem.Curriculumname;
				$scope.form.Curriculumsid = selectedItem.CurriculumsId;
			}else{
				$scope.form.CurriculumsName = "";
				$scope.form.Curriculumsid = null;
			}
			//
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'删除此课程计划将清除其下所有章节计划配置。<br/>您想继续吗？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        curriculumsclasscentredel(item.CurriculumsclasscentreId);
      }
    });
  }
    $scope.run = function(){
    	curriculumsclasscentrelist();
    }
    $scope.run();
}]);



/*    课程管理-课程计划-章节配置列表-课程配置弹窗      */
app.controller("modalKcpzContr", ['$scope', 'httpService','$modal', '$modalInstance','items','formValidate','toaster',function ($scope, httpService,$modal,$modalInstance,items,formValidate,toaster) {
    console.log("课程管理-课程计划-章节配置列表-课程配置弹窗")

    $scope.items = items;

    $scope.form = {
    	//  课程ID
    	"Curriculumsid":null,
    	//  课程名称
    	"Curriculumsname":"",
    	//  班级ID
    	"Classesid":null,
    	//  班级名称
    	"Classesname":"",
    	//  教师ID
    	"TeacherID":null,
    	//  教师名称
    	"Teachername":"",
    	//  是否录播[0:否,1:是]
    	"Isondemand":0,
    	// ----- 是否录播[0:否,1:是]
    	"IsondemandItem":{"val":0,"title":"否"},
    	// ----- 是否录播[0:否,1:是]
    	"IsondemandItems":[
    		{"val":0,"title":"否"},
    		{"val":1,"title":"是"}
    	],
    	//  是否直播[0:否,1:是]
    	"Islive":0,
    	// ------ 是否直播[0:否,1:是]
    	"IsliveItem":{"val":0,"title":"否"},
    	// ------ 是否直播[0:否,1:是]
    	"IsliveItems":[
    		{"val":0,"title":"否"},
    		{"val":1,"title":"是"}
    	]
    }

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


    //    课程班级配置添加
    var curriculumsclasscentreadd = function(){
		if(!(formValidate($scope.form.Curriculumsname).minLength(0).outMsg(2901).isOk))return false;
		if(!(formValidate($scope.form.Classesname).minLength(0).outMsg(2902).isOk))return false;
		if(!(formValidate($scope.form.Teachername).minLength(0).outMsg(2903).isOk))return false;
		var url=config.HttpUrl + "/system/us/curriculumsclasscentreadd";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Curriculumsid": Number($scope.form.Curriculumsid),
			"Classesid": Number($scope.form.Classesid),
			"TeacherID": Number($scope.form.TeacherID),
			"Isondemand": Number($scope.form.Isondemand),
			"Islive": Number($scope.form.Islive)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("课程班级配置添加",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '添加成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
    }


    //   课程班级配置修改
    var curriculumsclasscentrechange = function(id){
    	if(!id)return false;
    	var url=config.HttpUrl + "/system/us/curriculumsclasscentrechange";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Curriculumsid": Number($scope.form.Curriculumsid),
			"Classesid": Number($scope.form.Classesid),
			"TeacherID": Number($scope.form.TeacherID),
			"Isondemand": Number($scope.form.Isondemand),
			"Islive": Number($scope.form.Islive),
			"Id":Number(id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("课程班级配置修改",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '修改成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
    }


	//    打开弹窗  选择 课程选择 - 学科筛选课程
	$scope.modalOpenClassroom = function() {
		console.log("打开弹窗  选择 课程选择 - 学科筛选课程");
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_course_science.html',
			windowClass: 'modal-kcgl-kcjh-add',
			controller: 'modalGetCourseCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				},
				deps: ['$ocLazyLoad',
					function($ocLazyLoad) {
						return $ocLazyLoad.load(['ui.select']).then();
					}
				]
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!!selectedItem){
				$scope.form.Curriculumsname = selectedItem.Curriculumname;
				$scope.form.Curriculumsid = selectedItem.CurriculumsId;
			}else{
				$scope.form.Curriculumsname = "";
				$scope.form.Curriculumsid = null;
			}

		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}


	//    打开弹窗  选择班级
	$scope.modalOpenClass = function() {
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_class.html',
			controller: 'modalGetClassCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!!selectedItem){
				$scope.form.Classesname = selectedItem.Classesname;
				$scope.form.Classesid = selectedItem.Classid;
			}else{
				$scope.form.Classesname = "";
				$scope.form.Classesid = null;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}



	//  打开窗口-选择老师
	$scope.modalOpenTeacher=function(){
		console.log("打开窗口-选择老师");
		var modalInstance=$modal.open({
			templateUrl:'../html/modal/modal_teacher.html',
			controller:'modalGetTeacherCtrl',
			resolve:{
				items:function(){
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!!selectedItem){
				$scope.form.Teachername = selectedItem.Truename;
				$scope.form.TeacherID = selectedItem.Usersid;
			}else{
				$scope.form.Teachername = "";
				$scope.form.TeacherID = null;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}


	$scope.ok = function(){
		if(items.operate == "edit"){
			curriculumsclasscentrechange(items.item.CurriculumsclasscentreId);
		}else{
			curriculumsclasscentreadd();
		}
	}

	//取消按钮
	$scope.cancel=function(){
		$modalInstance.dismiss('cancel');
	}

	$scope.run = function(){
		if(items.operate == "edit"){
			$scope.form.Curriculumsid = items.item.Curriculumsid;
			$scope.form.Curriculumsname = items.item.Curriculumname;
			$scope.form.Classesid = items.item.Classesid;
			$scope.form.Classesname = items.item.Classesname;
			$scope.form.TeacherID = items.item.TeacherId;
			$scope.form.Teachername = items.item.Truename;
			$scope.form.Isondemand = items.item.Isondemand;
			$scope.form.Islive = items.item.Islive;
		}
	}
	$scope.run();


}]);
