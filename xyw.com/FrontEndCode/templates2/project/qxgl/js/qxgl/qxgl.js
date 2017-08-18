'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 权限管理
 */

/*    权限管理      */
app.controller("qxglContr", ['$scope','$state',function($scope, $state) {
	//$state.go("app.qxgl",false);
	console.log("权限管理")
	//   权限管理分页条数
	$scope.tableSize = 15;
}]);
