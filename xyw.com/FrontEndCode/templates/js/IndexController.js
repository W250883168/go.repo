/**
 * Created by Administrator on 2016/7/2.
 */



app.controller('indexController', function ($scope, httpService,notificService) {
    //服务器url
    $scope.serverUrl = "192.168.0.9:8090"

    //单个设备
    $scope.deviceCmd = {deviceId:"1",useWhoseCmd:"self",cmdCode:"on",para:""}
    $scope.deviceResponse = {code:"",msg:"",data:""}
    $scope.fnSendCmdToDevice = function () {
        var url = "http://" + $scope.serverUrl + "/device/control";
        var data = {
            DeviceId: $scope.deviceCmd.deviceId,         //设备id
            UseWhoseCmd: $scope.deviceCmd.useWhoseCmd,   //使用谁的命令：设备自己self/设备连接的节点node
                                        //象灯这种设备，本身是没有什么操作命令的，它的开关
                                        //控制是通过它连接的节点的开关操作来控制的，所以它
                                        //使用的是设备连接的节点的命令（node)。
            CmdCode: $scope.deviceCmd.cmdCode,           //命令代码
            Para: $scope.deviceCmd.para                 //参数,格式为："p1":"xxx","p2":"yyy"
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.deviceResponse = data;
        }, function (reason) {
        }, function (update) {
        })
    }

    //整个教室
    $scope.classroomCmd = {classroomId:"1",cmdCode:"on",para:""}
    $scope.classroomResponse = {code:"",msg:"",data:""}
    $scope.fnSendCmdToClassroom = function () {
        var url = "http://" + $scope.serverUrl + "/classroom/control";
        var data = {
            ClassroomId: $scope.classroomCmd.classroomId,    //教室id
            CmdCode: $scope.classroomCmd.cmdCode,           //命令代码
            Para: $scope.classroomCmd.para                 //参数,格式为："p1":"xxx","p2":"yyy"
        };
        var promise = httpService.ajaxPost(url, data);
        promise.then(function (data) {
            $scope.classroomResponse = data;
        }, function (reason) {
        }, function (update) {
        })
    }
})







