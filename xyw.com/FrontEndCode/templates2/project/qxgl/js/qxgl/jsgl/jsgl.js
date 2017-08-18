'use strict';
/**
 * Created by Administrator on 2016/7/28.
 *  权限管理-角色列表
 */
/*    角色列表     */
app.controller("qxglJsglContr", ['$scope','httpService','$modal','toaster', function($scope,httpService,$modal,toaster) {
  //$state.go("app.qxgl",false);
  console.log("角色列表")
  //   用户列表
  $scope.jsgl_list = [];
  //   page
  $scope.backPage = {
    PageIndex:1,
    PageSize:10
  }
  //   form
  $scope.form = {
    "KeyWord":"",
    "Id":0
  }
  //   教师列表
  var teacherlist = function(){
    var url = config.HttpUrl+"/system/sm/getroles";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "Rolesname":$scope.form.KeyWord,
      "Id":$scope.form.Id,
      "PageIndex": $scope.backPage.PageIndex,
      "PageSize": $scope.backPage.PageSize
    };
    var promise = httpService.ajaxPost(url,data);
    promise.then(function (data) {
      console.log("教师列表",data)
      if(data.Rcode=="1000"){
        $scope.jsgl_list = data.Result.PageData;
        //   分页
        var objPage={PageCount:data.Result.PageCount,PageIndex:data.Result.PageIndex,PageSize:data.Result.PageSize,RecordCount:data.Result.PageCount};
        if((objPage.RecordCount % objPage.PageSize)==0){
          objPage.PageCount=(objPage.RecordCount / objPage.PageSize);
        }else{
          objPage.PageCount=parseInt((objPage.RecordCount / objPage.PageSize))+1;
        }
        $scope.backPage = pageFn(objPage,5);
      }else{
      	$scope.jsgl_list = [];
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
    var url = config.HttpUrl + "/system/sm/delroles";
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
        toaster.pop('warning', data.Reason);
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
      $scope.searchPost();
    }
  }
  //添加按钮功能
  $scope.openModalAdd = function (str,item) {
    var modalInstance = $modal.open({
      templateUrl: '../project/qxgl/html/qxgl/jsgl/modal_add.html',
      controller: 'modalqxglContr',
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


/*    用户管理-弹窗     */
app.controller("modalqxglContr",['$scope', 'httpService', '$modalInstance','items','$modal','toaster',function ($scope, httpService,$modalInstance,items,$modal,toaster) {
  console.log("角色管理-弹窗");
  $scope.items = items;
  //   form
  $scope.form = {
    "Id":null,
    "Rolesname":""
  }

  //   查看  true==查看 。 false！=查看
  $scope.details = false;

  //    添加角色
  var usersadd = function(){
    var url = config.HttpUrl + "/system/sm/addroles";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "Rolesname": $scope.form.Rolesname
    };
    var promise = httpService.ajaxPost(url,data);
    promise.then(function (data) {
      console.log("添加角色",data)
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
    var url = config.HttpUrl + "/system/sm/updateroles";
    var data = {
      "Usersid": config.GetUser().Usersid,
      "Rolestype": config.GetUser().Rolestype,
      "Token": config.GetUser().Token,
      "Os": "WEB",
      "Rolesname": $scope.form.Rolesname,
      "Id": Number($scope.form.Id)
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

  $scope.cancel = function() {
    $modalInstance.dismiss('cancel');
  };
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
        $scope.form = $.extend({},$scope.form,$scope.items.item);
        break;
    }
  }
  $scope.run();
}]);
