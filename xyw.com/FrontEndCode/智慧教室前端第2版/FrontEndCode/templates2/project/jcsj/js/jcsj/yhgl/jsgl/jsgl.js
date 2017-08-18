'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 基础数据-用户管理-教师管理
 */

/*    教师管理     */
app.controller("jcsjYhglJsglContr", ['$scope','httpService','$modal','toaster', function($scope,httpService,$modal,toaster) {
	//$state.go("app.qxgl",false);
    console.log("教师管理")

    //   用户列表
    $scope.usersList = [];

    //   page
    $scope.backPage = {
    	PageIndex:1,
    	PageSize:10
    }

    //   form
    $scope.form = {
    	"KeyWord":"",
    	//   专业
    	"Majorids":"",
    	//   学院
    	"Collegeids":""
    }

    //   教师列表
    var teacherlist = function(){
        var url = config.HttpUrl+"/system/us/teacherlist";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Searhtxt":$scope.form.KeyWord,
            "Majorids":$scope.form.Majorids,
            "Collegeids":$scope.form.Collegeids,
			"PageIndex": $scope.backPage.PageIndex,
			"PageSize": $scope.backPage.PageSize
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("教师列表",data)
            if(data.Rcode=="1000"){
            	$scope.usersList = data.Result.PageData;
            	//   分页
				var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
            }else if(data.Rcode=="1002"){
            	$scope.usersList = [];
            	//   分页
				var objPage={PageCount:0,PageIndex:1,PageSize:10,RecordCount:0};
				if((objPage.RecordCount % objPage.PageSize)==0){
					objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
				}else{
					objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
				}
				$scope.backPage = pageFn(objPage,5);
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };


    //   删除用户
	var teacherdel = function(Id){
		var url = config.HttpUrl+"/system/us/teacherdel";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"Id": Id
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("删除用户",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '删除成功！');
            teacherlist();
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
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
			teacherlist();
		}
	};
	/*  -------------------- 分页、页码  -----------------------  */

    //   搜索
    $scope.searchPost = function(){
		$scope.backPage.PageIndex=1;
    	teacherlist();
    }

    //   回车查询
	$scope.sbgzKeyup = function(e){
        var keycode = window.event?e.keyCode:e.which;
        if(keycode==13){
            teacherlist();
        }
	}




	//添加按钮功能
    $scope.openModalAdd = function (str,item) {
        var modalInstance = $modal.open({
            templateUrl: '../project/jcsj/html/jcsj/yhgl/jsgl/modal_add.html',
            controller: 'modalYhglContr',
            windowClass: 'm-modal-yhgl',
            resolve: {
                items: function () {
                    return {'str':str,"item":item};
                }
            }
        });

        modalInstance.result.then(function(bul) {
			console.log(bul)
			if(bul){
				teacherlist();
			}

		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
    }


	//   点击删除
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
        teacherdel(item.UsersId);
      }
    });
  }

    $scope.run = function(){
    	teacherlist();
    }
    $scope.run();


}]);


/*    用户管理-弹窗     */
app.controller("modalYhglContr",['$scope', 'httpService', '$modalInstance','items','$modal','toaster','formValidate',function ($scope, httpService,$modalInstance,items,$modal,toaster,formValidate) {
	console.log("用户管理-弹窗");

	$scope.items = items;

	//   form
	$scope.form = {
		"UsersId":"",
		"Loginuser":"",
		"Loginpwd":"",
		"Rolesid":2,
		"RolesList":[],
		"Truename":"",
		"Nickname":"",
		"Userheadimg":"",
		"Userphone":"",
		"Userstate":null,
		"Usersex":1,
		"Birthday":"",
		"Collegeid":null,
		"Majorid":null,
		"CollegeMajorName":"",
		"Id":null,
		"RolesListV":"",
	    //  性别
	    "UsersexLists":[
	      {value:1,name:"女"},
	      {value:2,name:"男"}
	    ],
	    //  性别Id
	    UsersexList:""
	};

	//   查看  true==查看 。 false！=查看
	$scope.details = false;

	//   上传图片
	$scope.upimglist = [];

	//    清除图片
	$scope.closePic = function(index){
		$scope.upimglist.splice(index,1);
	}
	//   用户角色列表
	var getroles = function(){
        var url = config.HttpUrl+"/system/sm/getroles";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
			"PageIndex": -1
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("用户角色列表",data)
            if(data.Rcode=="1000"){
                if(data.Result.PageData != null)$scope.form.RolesList = data.Result.PageData;
            }else{
              toaster.pop('warning',data.Reason);
            }
        }, function (reason) {}, function (update) {});
    };

	//    添加用户
	var usersadd = function(){
		if(!(formValidate($scope.form.Loginuser).minLength(0).outMsg(3101).isOk)) return false;
		if(!(formValidate($scope.form.Loginpwd).minLength(0).outMsg(3102).isOk)) return false;
		if(!(formValidate($scope.form.RolesListV.Id).isNumber().outMsg(3103).isOk)) return false;
		if(!(formValidate($scope.form.Truename).minLength(0).outMsg(3104).isOk)) return false;
		if(!(formValidate($scope.form.CollegeMajorName).minLength(0).outMsg(3105).isOk)) return false;
		if(!(formValidate($scope.form.Nickname).minLength(0).outMsg(3106).isOk)) return false;
		if(!(formValidate($scope.form.Userphone).isMobile().outMsg(2205).isOk)) return false;
        var url = config.HttpUrl+"/system/us/teacheradd";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Rolesid":Number($scope.form.Rolesid),
            "Loginuser":jQuery.trim($scope.form.Loginuser),
            "Loginpwd":jQuery.trim($scope.form.Loginpwd),
            "Truename":$scope.form.Truename,
            "Nickname":$scope.form.Nickname,
            "Userheadimg":$scope.form.Userheadimg,
            "Userphone":$scope.form.Userphone,
            "Userstate":$scope.form.Userstate,
            "Usersex":Number($scope.form.Usersex),
            "Birthday":$scope.form.Birthday,
            "Collegeid":Number($scope.form.Collegeid),
			"Majorid":Number($scope.form.Majorid)
		};
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("添加用户",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '添加成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };

	//    修改用户
	var userschange = function(){
		if(!(formValidate($scope.form.Loginuser).minLength(0).outMsg(3101).isOk)) return false;
		if(!(formValidate($scope.form.RolesListV.Id).isNumber().outMsg(3103).isOk)) return false;
		if(!(formValidate($scope.form.Truename).minLength(0).outMsg(3104).isOk)) return false;
		if(!(formValidate($scope.form.CollegeMajorName).minLength(0).outMsg(3105).isOk)) return false;
		if(!(formValidate($scope.form.Nickname).minLength(0).outMsg(3106).isOk)) return false;
		if(!(formValidate($scope.form.Userphone).isMobile().outMsg(2205).isOk)) return false;
        var url = config.HttpUrl+"/system/us/teacherchange";
        var data = {
            "Usersid": config.GetUser().Usersid,
            "Rolestype": config.GetUser().Rolestype,
            "Token": config.GetUser().Token,
            "Os": "WEB",
            "Rolesid":Number($scope.form.Rolesid),
            "Loginuser":jQuery.trim($scope.form.Loginuser),
            "Loginpwd":jQuery.trim($scope.form.Loginpwd),
            "Truename":$scope.form.Truename,
            "Nickname":$scope.form.Nickname,
            "Userheadimg":$scope.form.Userheadimg,
            "Userphone":$scope.form.Userphone,
            "Userstate":$scope.form.Userstate,
            "Usersex":Number($scope.form.Usersex),
            "Birthday":$scope.form.Birthday,
            "Collegeid":Number($scope.form.Collegeid),
	      	"Majorid":Number($scope.form.Majorid),
            "Id":$scope.form.UsersId
		    };
        var promise = httpService.ajaxPost(url,data);
        promise.then(function (data) {
        	console.log("修改用户",data)
          if(data.Rcode == "1000"){
            toaster.pop('success', '修改成功！');
            $modalInstance.close(true);
          }else{
            toaster.pop('warning',data.Reason);
          }
        }, function (reason) {}, function (update) {});
    };


	//    选择时间
	$scope.showDate = function(){
		jeDate({
			dateCell: "#jcsj_yhgl_data",
			format: "YYYY-MM-DD",
			isTime: true,
			minDate: "1900-12-31 00:00:00",
			isinitVal: false,
			choosefun:function(elem,val) {
				$scope.form.Birthday = val;
			},
			okfun: function(elem,val) {
				$scope.form.Birthday = val;
			},
			clearfun:function(elem, val) {
				$scope.form.Birthday = "";
			}
		});
	}

	//    打开修改密码
	$scope.openModalPwd = function(){
		var modalInstance = $modal.open({
            templateUrl: '../project/jcsj/html/jcsj/yhgl/modal_add_pwd.html',
            controller: 'modalYhglPwdContr',
            windowClass: 'm-modal-yhgl',
            resolve: {
                items: function () {
                    return $scope.items.item;
                }
            }
        });

        modalInstance.result.then(function(bul) {
			console.log(bul)

		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}


	//    打开弹窗  选择班级 / 学科
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
				$scope.form.Collegeid = selectedItem.Collegeid;
				$scope.form.Majorid = selectedItem.Majorid;
				if(selectedItem.Collegename.length > 0){
					$scope.form.CollegeMajorName = selectedItem.Collegename;
					if(selectedItem.Majorname.length > 0)$scope.form.CollegeMajorName += "-" + selectedItem.Majorname;
				}else{
					$scope.form.CollegeMajorName = "";
				}

			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	}


	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	};
	// 用户角色选择传ID
	$scope.changeModelItem = function(item) {
		console.log('用户角色选择传ID', item);
		$scope.form.Rolesid = item.Id;
	};
  // 性别选择传ID
  $scope.changeUsersexList = function (item) {
    $scope.form.Usersex = item.value;
  }
	//   ok
	$scope.ok = function(){
		switch($scope.items.str){
			case "add":
				if($scope.upimglist.length > 0){
					$scope.form.Userheadimg = $scope.upimglist[0].Result;
				}else{
					$scope.form.Userheadimg = "";
				}
				usersadd();
			break;
			case "edit":
				if($scope.upimglist.length > 0){
					$scope.form.Userheadimg = $scope.upimglist[0].Result;
				}else{
					$scope.form.Userheadimg = "";
				}
				userschange();
			break;
		}
	}



	$scope.run = function(){
		//getroles();
		//   用户角色列表
		$scope.form.RolesList = [{
			"Id":2,
			"Rolesname":"教师人员"
		}];

		switch($scope.items.str){
			case "add":

			break;
			case "details":
				$scope.details = true;

				$scope.form = $.extend({},$scope.form,$scope.items.item);
				//    属所院系
				if($scope.form.CollegeName.length > 0){
					$scope.form.CollegeMajorName = $scope.form.CollegeName;
					if($scope.form.MajorName.length > 0)$scope.form.CollegeMajorName += "-" + $scope.form.MajorName;
				}else{
					$scope.form.CollegeMajorName = $scope.form.MajorName;
				}
				if($scope.form.Userheadimg.length > 0)$scope.upimglist[0] = {Result:$scope.form.Userheadimg};
        //  选择性别 当前的状态
        for(var item in $scope.form.UsersexLists){
          if($scope.form.UsersexLists[item].value == $scope.form.Usersex) {
            $scope.form.UsersexList = $scope.form.UsersexLists[item]
          }
        }
        //   选择用户角色 当前的状态
        for (var item in $scope.form.RolesList) {
          if ($scope.form.RolesList[item].Id == $scope.form.Rolesid) {
            $scope.form.RolesListV = $scope.form.RolesList[item];
          }
        }
			break;
			case "edit":
				$scope.details = false;
				$scope.form = $.extend({},$scope.form,$scope.items.item);
				//    属所院系
				if($scope.form.CollegeName.length > 0){
					$scope.form.CollegeMajorName = $scope.form.CollegeName;
					if($scope.form.MajorName.length > 0)$scope.form.CollegeMajorName += "-" + $scope.form.MajorName;
				}else{
					$scope.form.CollegeMajorName = $scope.form.MajorName;
				}

				if($scope.form.Userheadimg.length > 0)$scope.upimglist[0] = {Result:$scope.form.Userheadimg};
        //  选择性别 当前的状态
        for(var item in $scope.form.UsersexLists){
          if($scope.form.UsersexLists[item].value == $scope.form.Usersex) {
            $scope.form.UsersexList = $scope.form.UsersexLists[item]
          }
        }
        //   选择用户角色 当前的状态
        for (var item in $scope.form.RolesList) {
          if ($scope.form.RolesList[item].Id == $scope.form.Rolesid) {
            $scope.form.RolesListV = $scope.form.RolesList[item];
          }
        }
			break;
		}
	}
	$scope.run();

}]);



