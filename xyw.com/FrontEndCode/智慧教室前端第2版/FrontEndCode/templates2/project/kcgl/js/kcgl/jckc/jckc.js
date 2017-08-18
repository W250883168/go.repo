'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 课程管理
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

/*    课程管理-基础课程      */
app.controller("kcglJckcContr", ['$scope', 'httpService', '$modal','$state', 'toaster',function ($scope, httpService, $modal,$state,toaster) {
    console.log("课程管理-基础课程")

    //   课程名称模板搜索 查询关键词
    $scope.form = {
    	Curriculumname:""
    }
    //   课程性质
    $scope.Curriculumnature = "";
    //   课程类型
    $scope.Curriculumstype = "";
    //   基础课程列表数组
    $scope.Curriculumlist = [];
    //   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:10
    }

    //   取课程列表
    var curriculumslist = function(){
    	var url=config.HttpUrl+"/system/us/curriculumslist";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Curriculumname": $scope.form.Curriculumname,
			"Curriculumnature": $scope.Curriculumnature,
			"Curriculumstype": $scope.Curriculumstype,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取课程列表",data)
			if(data.Rcode == "1000") {
				$scope.Curriculumlist = data.Result.PageData;
				//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
			}else if(data.Rcode=="1002"){
            	$scope.Curriculumlist = [];
            	//   分页
				var objPage={PageCount:0,PageIndex:$scope.backPage.PageIndex,PageSize:$scope.backPage.PageSize,RecordCount:0};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
            }else{
              toaster.pop('warning',data.Reason);
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
			curriculumslist();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */


    //    查询
    $scope.searchPost = function(){
		$scope.backPage.PageIndex=1;
    	curriculumslist();
    }
    //   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            curriculumslist();
        }
	}

	//   删除课程列表
	var itemDelete = function(Id){
		var url=config.HttpUrl+"/system/us/curriculumsdel";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Number(Id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除课程",data)
			if(data.Rcode == "1000") {
				//   page
			    $scope.backPage = {
			    	PageIndex:1,
			    	PageSize:10
			    }
				curriculumslist();
        toaster.pop("success","删除成功！");
			} else {
        toaster.pop("warning",data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}
  //   删除课程弹窗
  $scope.deleteItem = function(item){
    var modalInstance = $modal.open({
      templateUrl: 'modal/modal_alert_all.html',
      controller: 'modalAlert2Conter',
      resolve: {
        items: function () {
          return {"type":'warning',"msg":'删除此课程将清除其下所有章节及课程计划配置<br/>您想继续吗？'};
        }
      }
    });
    modalInstance.result.then(function(bul){
      if(bul){
        itemDelete(item.CurriculumsId);
      }
    });
  }

    $scope.run = function(){
    	curriculumslist();
    }
    $scope.run();


}]);


/*    课程管理-基础课程-添加课程      */
app.controller("modalJckcAddContr", ['$scope', 'httpService','$modal','$location','$state','toaster','formValidate', function ($scope, httpService,$modal,$location,$state,toaster,formValidate) {
    console.log("课程管理-基础课程-添加课程")

    $scope.form = {
    	//   课程ID
    	"Curriculumsid":null,
		//   课程
		"Curriculumsinfo":{},
		//   章节
		"Chapterslist":[]
	}
	$scope.form.Curriculumsinfo = {
		"Curriculumname":"",
		"Curriculumicon":"",
		"Curriculumnature":"",
		"Curriculumstype":"普通课",
		//  ---
		"CurriculumstypeItem":"",
		//  ---
		"CurriculumstypeItems":[
			{"val":"1","title":"普通课"},
			{"val":"2","title":"公开课"}
		],
		"Curriculumsdetails":"",
		"Chaptercount":"",
		"Subjectcode":"",
		"Subjectname":"",
		"Averageclassrate":"",
		"Createdate":""
	}
	//   章节
	var Chapters = {
		"Chaptername":"",
		"ChaptersIndex":50,
		"Curriculumsid":0,
		"Chaptericon":"",
		"Chapterdetails":""//,
		//"Createdate":""
	}

	//   操作
    var operation = $location.search().op;
    $scope.operation = $location.search().op;
    //   课程ID
    $scope.form.Curriculumsid = $location.search().cid;

	//   上传图片
	$scope.upimglist = [];

	//   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:10
    }

	//    添加 查看 修改
	$scope.purview = {add:true,details:true,edit:true};

	//    课程类型
	$scope.changeCurriculumstypeItem = function(item){
		$scope.form.Curriculumsinfo.CurriculumstypeItem = item;
		$scope.form.Curriculumsinfo.Curriculumstype = item.title;
	}

    //    打开 添加章节信息弹窗
    $scope.openModalAddZj = function (str,item) {
        var modalInstance = $modal.open({
            templateUrl: '../project/kcgl/html/kcgl/jckc/modal_add_add.html',
            controller: 'modalJckcAddZjContr',
            windowClass: 'm-modal-kcgl-zjpz',
            resolve: {
                items: function () {
                    return {form:$scope.form,operation:str,item:item};
                }
            }
        });
        modalInstance.result.then(function(item2) {
			if(item2 === true){
				if(operation == "edit")chapterslist();
			}else{
				if(operation == "add" && str == "edit"){
					item.Chapterdetails = item2.Chapterdetails;
					item.Chaptericon = item2.Chaptericon;
					item.Chaptername = item2.Chaptername;
					item.ChaptersIndex = item2.ChaptersIndex;
					item.Curriculumsid = item2.Curriculumsid;
				}else{
					var temp = $.extend({}, Chapters, item2);
					//    add
					$scope.form.Chapterslist.push(temp);
				}
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
    }

	//    打开弹窗  选择学科
	$scope.modalOpenClassroom = function() {
		console.log("打开弹窗  选择学科");
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_science.html',
			controller: 'modalGetScienceCtrl',
			resolve: {
				items: function() {
					return $scope.items;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem){
				//$scope.addAdd = "";
			}else{
				$scope.form.Curriculumsinfo.Subjectcode = selectedItem.Subjectcode;
				$scope.form.Curriculumsinfo.Subjectname = selectedItem.Subjectname;
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}



	/*  ------- add ------  */

	//    课程添加
	var curriculumsadd = function(){

		if(!(formValidate($scope.form.Curriculumsinfo.Subjectcode).minLength(0).outMsg('所属学科不能为空！').isOk))return false;
		if(!(formValidate($scope.form.Curriculumsinfo.Curriculumname).minLength(0).outMsg('课程名称不能为空！').isOk))return false;
		if(!(formValidate($scope.form.Curriculumsinfo.Curriculumstype).minLength(0).outMsg('课程类型不能为空！').isOk))return false;

		//   有图片
		if($scope.upimglist.length > 0){
			$scope.form.Curriculumsinfo.Curriculumicon = $scope.upimglist[0].Result;
		}else{
			$scope.form.Curriculumsinfo.Curriculumicon = "";
		}

		var url=config.HttpUrl+"/system/us/curriculumsadd";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Curriculumsinfo":{
            	"Curriculumname":$scope.form.Curriculumsinfo.Curriculumname,
            	"Curriculumicon":$scope.form.Curriculumsinfo.Curriculumicon,
            	"Curriculumnature":$scope.form.Curriculumsinfo.Curriculumnature,
            	"Curriculumstype":$scope.form.Curriculumsinfo.Curriculumstype,
            	"Curriculumsdetails":$scope.form.Curriculumsinfo.Curriculumsdetails,
            	"Subjectcode":$scope.form.Curriculumsinfo.Subjectcode
            },
            "Chapterslist":$scope.form.Chapterslist
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取课程列表",data)
			if(data.Rcode == "1000") {
				//$scope.Curriculumlist = data.Result.PageData;
				//   返回
        toaster.pop('success',"添加成功！");
				$state.go("app.kcgl.jckc",{},{ reload: true });
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}



	/*  ------- details ------  */
	//   取单条课程
    var curriculumsinfo = function(){
    	var url=config.HttpUrl+"/system/us/curriculumsinfo";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Curriculumsid": Number($scope.form.Curriculumsid)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取单条课程",data)
			if(data.Rcode == "1000") {
				$scope.form.Curriculumsinfo = data.Result.PageData;
				//   载入查看时图片
				if($scope.form.Curriculumsinfo.Curriculumicon != "")$scope.upimglist.push({Result:$scope.form.Curriculumsinfo.Curriculumicon});
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
    }

    //   取章节列表
    var chapterslist = function(){
    	var url=config.HttpUrl+"/system/us/chapterslist";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Curriculumsid": Number($scope.form.Curriculumsid),
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取章节列表",data)
			if(data.Rcode == "1000") {
				$scope.form.Chapterslist = data.Result.PageData;
				//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
			} else {
				$scope.form.Chapterslist = [];
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
    }

    /*  ------- edit ------  */
   //    课程修改
	var curriculumschange = function(){

		if(!(formValidate($scope.form.Curriculumsinfo.Subjectcode).minLength(0).outMsg('所属学科不能为空！').isOk))return false;
		if(!(formValidate($scope.form.Curriculumsinfo.Curriculumname).minLength(0).outMsg('课程名称不能为空！').isOk))return false;
		if(!(formValidate($scope.form.Curriculumsinfo.Curriculumstype).minLength(0).outMsg('课程类型不能为空！').isOk))return false;

		//   有图片
		if($scope.upimglist.length > 0){
			$scope.form.Curriculumsinfo.Curriculumicon = $scope.upimglist[0].Result;
		}else{
			$scope.form.Curriculumsinfo.Curriculumicon = "";
		}

		var url=config.HttpUrl+"/system/us/curriculumschange";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Curriculumname":$scope.form.Curriculumsinfo.Curriculumname,
            "Curriculumicon":$scope.form.Curriculumsinfo.Curriculumicon,
            "Curriculumnature":$scope.form.Curriculumsinfo.Curriculumnature,
            "Curriculumstype":$scope.form.Curriculumsinfo.Curriculumstype,
            "Curriculumsdetails":$scope.form.Curriculumsinfo.Curriculumsdetails,
            "Chaptercount":$scope.form.Curriculumsinfo.Chaptercount,
            "Subjectcode":$scope.form.Curriculumsinfo.Subjectcode,
            "Averageclassrate":$scope.form.Curriculumsinfo.Averageclassrate,
            "Createdate":$scope.form.Curriculumsinfo.Createdate,
            "Id":Number($scope.form.Curriculumsid)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("课程修改",data)
			if(data.Rcode == "1000") {
				//$scope.Curriculumlist = data.Result.PageData;
				//   返回
        toaster.pop('success',"修改成功！");
				$state.go("app.kcgl.jckc",{},{ reload: true });
			} else {
        toaster.pop('warning',data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}




    /*  ------- end ------  */

    //    章节操作
    $scope.operationClick = function(str,item){
    	$scope.openModalAddZj(str,item);
    }

	//    章节删除
	var chaptersdel = function(id){
		var url=config.HttpUrl+"/system/us/chaptersdel";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Number(id)
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("章节删除",data)
			if(data.Rcode == "1000") {
				chapterslist();
        toaster.pop('success',"删除成功！");
      } else {
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
	}


	//    删除图片
	$scope.closePic = function(index){
		$scope.upimglist.splice(index,1);
	}


	$scope.ok = function(){

		//    操作
		switch(operation){
			case "add":
				curriculumsadd();
			break;
			case "details":
			break;
			case "edit":
				curriculumschange();
			break;
		}
	}

	//	取消按钮
	$scope.cancel = function(){
		$state.go("app.kcgl.jckc");
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
			chapterslist();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */


	/**
	 * 删除
	 */
  $scope.chapterDelete = function(item,index){
    if(operation == "add"){
      var modalInstance = $modal.open({
        templateUrl: 'modal/modal_alert_all.html',
        controller: 'modalAlert2Conter',
        resolve: {
          items: function () {
            return {"type":'warning',"msg":'你确定要删除临时章节吗？'};
          }
        }
      });
      modalInstance.result.then(function(bul){
        if(bul){
          $scope.form.Chapterslist.splice(index, 1);
        }
      });
    } else {
      var modalInstance = $modal.open({
        templateUrl: 'modal/modal_alert_all.html',
        controller: 'modalAlert2Conter',
        resolve: {
          items: function () {
            return {"type":'warning',"msg":'删除此章节将清除课程计划中章节配置<br/>您想继续吗？'};
          }
        }
      });
      modalInstance.result.then(function(bul){
        if(bul){
          chaptersdel(item.Id);
        }
      });
    }
  }

	$scope.run = function(){
		//   课程类型
		$scope.form.Curriculumsinfo.CurriculumstypeItem = $scope.form.Curriculumsinfo.CurriculumstypeItems[0];
		$scope.form.Curriculumsinfo.Curriculumstype = $scope.form.Curriculumsinfo.CurriculumstypeItems[0].title;

		//    操作
		switch(operation){
			case "add":
				$scope.purview = {add:true,details:false,edit:false};
				//   面包屑标题
				$scope.operatetitle = "添加课程";
			break;
			case "details":
				$scope.purview = {add:false,details:true,edit:false};
				curriculumsinfo();
				chapterslist();
				//   面包屑标题
				$scope.operatetitle = "查看课程";
			break;
			case "edit":
				$scope.purview = {add:false,details:false,edit:true};
				curriculumsinfo();
				chapterslist();
				//   面包屑标题
				$scope.operatetitle = "修改课程";
			break;
		}
	}
	$scope.run();

}]);


/*    课程管理-基础课程-添加课程-添加章节弹窗      */
app.controller("modalJckcAddZjContr", ['$scope', 'httpService', '$modal','items','$modalInstance','$location','formValidate', 'toaster',function ($scope, httpService, $modal,items,$modalInstance,$location,formValidate,toaster) {
    console.log("课程管理-课程计划-添加课程-添加章节弹窗")
    //   操作
    var operation = $location.search().op;
    console.log(items);

    //   章节
	$scope.Chapters = {
		"Chaptername":"",
		"ChaptersIndex":50,
		"Curriculumsid":0,
		"Chaptericon":"",
		"Chapterdetails":""//,
		//"Createdate":""
	}

	//    课程
    $scope.Curriculumsinfo = items.form.Curriculumsinfo;
    //   item
    $scope.item = items.item;
    $scope.items = items;

	//   上传图片
	$scope.upimglist = [];

	//    添加 查看 修改
	$scope.purview = {add:true,details:true,edit:true};


	//    单条章节添加
	var chaptersadd = function(item){

		if(!item && item.Curriculumsid != null)item.Curriculumsid = Number(item.Curriculumsid);
		item.Curriculumsid = Number($location.search().cid);
    	var url=config.HttpUrl+"/system/us/chaptersadd";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB"
		};
		data = $.extend({},data,item);
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("单条章节添加",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '添加成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
    }

	//    单条章节修改
	var chapterschange = function(item){
		if(!item && item.Curriculumsid != null)item.Curriculumsid = Number(item.Curriculumsid);
		item.Curriculumsid = Number($location.search().cid);
		item.Id = items.item.Id;
    	var url=config.HttpUrl+"/system/us/chapterschange";
		var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB"
		};
		data = $.extend({},data,item);
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("单条章节修改",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '修改成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
		}, function(reason) {}, function(update) {});
    }


	$scope.ok = function () {
		if(!(formValidate($scope.Chapters.Chaptername).minLength(0).outMsg(2900).isOk))return false;
		//
		if($scope.upimglist.length > 0){
			//   传入图片
			$scope.Chapters.Chaptericon = $scope.upimglist[0].Result;
		}else{
			$scope.Chapters.Chaptericon = "";
		}
		//   排序
		if(!!$scope.Chapters.ChaptersIndex)$scope.Chapters.ChaptersIndex = Number($scope.Chapters.ChaptersIndex);
		//   课程 不等于添加
		if(operation != "add"){
			if(items.operation == "add")chaptersadd($scope.Chapters);
			if(items.operation == "edit"){

				//    章节修改
				chapterschange($scope.Chapters);
			}
		}else{
			$modalInstance.close($scope.Chapters);
		}

    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };

    //    删除图片
	$scope.closePic = function(index){
		$scope.upimglist.splice(index,1);
	}

	$scope.run = function(){
		switch(items.operation){
			case "add":
				//    添加 查看 修改
				$scope.purview = {add:true,details:false,edit:false};
				//   序号
				$scope.Chapters.ChaptersIndex = $scope.items.form.Chapterslist.length + 1;
			break;
			case "details":
				$scope.Chapters = $.extend({},$scope.Chapters,$scope.item);
				//   载入查看时图片
				$scope.upimglist.push({Result:$scope.Chapters.Chaptericon});
				$scope.purview = {add:false,details:true,edit:false};
			break;
			case "edit":
				$scope.Chapters = $.extend({},$scope.Chapters,$scope.item);
				//   载入查看时图片
				if($scope.Chapters.Chaptericon)$scope.upimglist.push({Result:$scope.Chapters.Chaptericon});
				$scope.purview = {add:false,details:false,edit:true};
			break;
		}
	}
	$scope.run();
}]);
