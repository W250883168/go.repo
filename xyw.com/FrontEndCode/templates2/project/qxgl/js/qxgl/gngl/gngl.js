'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 功能管理
 */
/*  ////////////////////////////////////////////////////////////  */

/*  //////////////////////////////////////////////////  */
/*    功能管理      */
app.controller("qxglGnglContr", ['$scope','httpService','$modal','toaster', function($scope,httpService,$modal,toaster) {
  //   模块管理
  $scope.gngl_list = [];
  //   page
  $scope.backPage = {
    PageIndex:1,
    PageSize:10
  }
  //   form
  $scope.form = {
    "Functionname":"",
    "Id":0,
    "Superiormoduleid":null
  }
  //   模块列表
  var teacherlist = function(){
    var url = config.HttpUrl + "/system/sm/getsystemmodelfunc";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "Functionname":$scope.form.Functionname,
      "Id":Number($scope.form.Id),
      "Superiormoduleid":Number($scope.form.Superiormoduleid),
      "PageIndex": $scope.backPage.PageIndex,
      "PageSize": $scope.backPage.PageSize
    };
    var promise = httpService.ajaxPost(url,data);
    promise.then(function (data) {
      if(data.Rcode=="1000"){
        $scope.gngl_list = data.Result.PageData;
        //   分页
        var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
        if((objPage.RecordCount % objPage.PageSize)==0){
          objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
        }else{
          objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
        }
        $scope.backPage = pageFn(objPage,5);
      }else{
      	$scope.gngl_list = [];
        //   分页
        var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
        if((objPage.RecordCount % objPage.PageSize)==0){
          objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
        }else{
          objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
        }
        $scope.backPage = pageFn(objPage,5);
        //toaster.pop('warning',data.Reason);
      }
    }, function (reason) {}, function (update) {});
  };


  //   删除用户
  var teacherdel = function(Id){
    var url = config.HttpUrl + "/system/sm/delsystemmodelfunc";
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
      if(data.Rcode=="1000"){
        teacherlist();
        toaster.pop('success', '删除成功！');
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
  	$scope.backPage.PageIndex = 1;
    teacherlist();
  }

  //   回车查询
  $scope.sbgzKeyup = function(e){
    var keycode = window.event?e.keyCode:e.which;
    if(keycode == 13){
      $scope.searchPost();
    }
  }
  //添加按钮功能
  $scope.openModalAdd = function (str,item) {
    var modalInstance = $modal.open({
      templateUrl: '../project/qxgl/html/qxgl/gngl/modal_add.html',
      controller: 'modalgnglContr',
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
  //  删除弹窗功能
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
        $scope.backPage.PageIndex = 1;
        $scope.backPage.PageSize = 10;
        teacherdel(item.Id);
      }
    });
  }
  $scope.run = function(){
    teacherlist();
  }
  $scope.run();
}]);

/*    功能管理-弹窗      */
app.controller("modalgnglContr",['$scope', 'httpService', '$modalInstance','items','$modal','toaster',function ($scope, httpService,$modalInstance,items,$modal,toaster) {
  $scope.items = items;
  console.log("功能管理-弹窗",items);
  //   form
  $scope.form = {
    "Systemmoduleid":null,
    "Id":null,
    "Functionicon":"",
    "Functionname":"",
    "Functioncode":"",
    "mkglLists": [],
    "Functionsurls":"",
    "Functionsattribute":"",
    "FunctionDescribe":""
  }

  //   查看  true==查看 。 false！=查看
  $scope.details = false;

  //    添加功能
  var usersadd = function(){
    var url = config.HttpUrl + "/system/sm/addsystemmodelfunc";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "Id":Number($scope.form.Id),
      "Systemmoduleid":Number($scope.form.Systemmoduleid),
      "Functionicon": $scope.form.Functionicon,
      "Functionname": $scope.form.Functionname,
      "Functioncode": $scope.form.Functioncode,
      "Functionsurls": $scope.form.Functionsurls,
      "Functionsattribute": $scope.form.Functionsattribute,
      "FunctionDescribe": $scope.form.FunctionDescribe
    };
    var promise = httpService.ajaxPost(url,data);
    promise.then(function (data) {
      console.log("添加模块",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '添加成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
    }, function (reason) {}, function (update) {});
  };

  //    修改功能
  var userschange = function(){
    var url = config.HttpUrl + "/system/sm/updatesystemmodelfunc";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "Id":Number($scope.form.Id),
      "Systemmoduleid":Number($scope.form.Systemmoduleid),
      "Functionicon": $scope.form.Functionicon,
      "Functionname": $scope.form.Functionname,
      "Functioncode": $scope.form.Functioncode,
      "Functionsurls": $scope.form.Functionsurls,
      "Functionsattribute": $scope.form.Functionsattribute,
      "FunctionDescribe": $scope.form.FunctionDescribe
    };
    var promise = httpService.ajaxPost(url,data);
    promise.then(function (data) {
      console.log("修改修改功能",data)
      if(data.Rcode == "1000"){
        toaster.pop('success', '修改成功！');
        $modalInstance.close(true);
      }else{
        toaster.pop('warning',data.Reason);
      }
    }, function (reason) {}, function (update) {});
  };
  //  取模块列表
  var teacherlist = function(){
    var url = config.HttpUrl + "/system/sm/getsystemmodel";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB"
    };
    var promise = httpService.ajaxPost(url,data);
    promise.then(function (data) {
      if(data.Rcode=="1000"){
        $scope.form.mkglLists = data.Result.PageData;

      }else if(data.Rcode=="1002"){
        $scope.form.mkglLists = [];
      }else{
        toaster.pop('warning',data.Reason);
      }
    }, function (reason) {}, function (update) {});
  };

  $scope.cancel = function() {
    $modalInstance.dismiss('cancel');
  };
  $scope.changeModelItem = function (item) {
    $scope.form.Systemmoduleid = item.id;
  }
  //   ok
  $scope.ok = function(){
    switch($scope.items.str){
      case "add":
        usersadd();
        break;
      case "edit":
        userschange();
        break;
    }
  }

  $scope.run = function(){
    switch($scope.items.str){
      case "add":

        break;
      case "details":
        $scope.details = true;
        $scope.form = $.extend({},$scope.form,$scope.items.item);
        break;
      case "edit":
        $scope.details = false;
        $scope.form.Systemmoduleid = $scope.items.item.Systemmoduleid;
        $scope.form = $.extend({},$scope.form,$scope.items.item);
        break;
    }
    teacherlist();
  }
  $scope.run();
}]);
