app.factory('notificService', function () {
    return {
        //type: teal amethyst ruby tangerine lemon  lime ebony smoke
        notific: function (type, sticky, title, msg,time) {
            var settings = {
                theme: type,                //样式
                sticky: sticky,             //是否黏住：true:是（不自动关闭），false:否（自动关闭）
                horizontalEdge: 'top',      //位置：top bottom 
                verticalEdge: 'righ',       //位置：left right  
                heading: title,             //标题
                life: time                  //显示时间:毫秒
            }

            $.notific8('zindex', 11500);
            $.notific8(msg, settings);

        }

    };
});