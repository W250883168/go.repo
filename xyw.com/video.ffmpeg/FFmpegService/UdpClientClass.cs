using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;
using System.Timers;


namespace FFmpegService
{
    /// <summary>
    /// ffmpeg服务器
    /// </summary>
    public class UdpClientClass : IDisposable
    {
        FFmpegCmd cmd = new FFmpegCmd();

        /// <summary>  
        /// 构建客户端  
        /// </summary>  
        /// <param name="server">服务器iP地址或者域名</param>  
        /// <param name="port">服务器监听端口</param>  
        /// <param name="locadPort">本地监听端口</param>  
        /// <param name="timeout">超时等待时间</param>  
        public UdpClientClass(string server, int port, int locadPort, int timeout)
        {
            try
            {
                ServerIPE = new IPEndPoint(IPAddress.Parse(server), port);
                UdpListenClient = new UdpClient(locadPort);//固定通信端口  
                UdpListenClient.Client.ReceiveTimeout = timeout;
                const long IOC_IN = 0x80000000;
                const long IOC_VENDOR = 0x18000000;
                const long SIO_UDP_CONNRESET = IOC_IN | IOC_VENDOR | 12;
                byte[] optionInValue = { Convert.ToByte(false) };
                byte[] optionOutValue = new byte[4];
                UdpListenClient.Client.IOControl((IOControlCode)SIO_UDP_CONNRESET, optionInValue, optionOutValue);
            }
            catch (System.Exception ex)
            {
                string dir = System.Windows.Forms.Application.StartupPath;
                string file = "log_" + DateTime.Now.ToString("yyyyMMdd") + ".log";
                string path = System.IO.Path.Combine(dir, file);
                using (System.IO.StreamWriter sw = new System.IO.StreamWriter(path, true))
                {
                    sw.WriteLine(ex.Message.ToString());
                    sw.WriteLine(ex.InnerException);
                }
            }
        }

        /// <summary>  
        /// UDP发送类，绑定了一个固定的端口  
        /// </summary>  
        private static UdpClient UdpListenClient;

        /// <summary>  
        /// 服务器端的IP与端口  
        /// </summary>  
        private IPEndPoint ServerIPE = null;
        bool IsReceiving = false;

        public void Send(byte[] data, int len)
        {
            int sends = UdpListenClient.Send(data, len, ServerIPE);
            if (!IsReceiving)
                StartAndLsn();
        }

        private Thread ClientRecThread;

        private void StartAndLsn()
        {
            IsReceiving = true;
            ClientRecThread = new Thread(new ThreadStart(ThreadFunc_Recv));//启动新线程做接收  
            ClientRecThread.IsBackground = true;
            ClientRecThread.Start();
        }

        private void ThreadFunc_Recv()//接收数据做服务  
        {
            string dir = System.Windows.Forms.Application.StartupPath;
            string file = "log_" + DateTime.Now.ToString("yyyyMMdd") + ".log";
            string path = System.IO.Path.Combine(dir, file);
            using (System.IO.StreamWriter sw = new System.IO.StreamWriter(path, true))
            {
                while (IsReceiving)
                {
                    IPAddress ipaddr = ServerIPE.Address;
                    IPEndPoint remoteIPE = new IPEndPoint(ipaddr, 0);
                    int buffSizeCurrent = UdpListenClient.Client.Available;
                    if (buffSizeCurrent > 0)
                    {
                        try
                        {
                            byte[] data_recv = UdpListenClient.Receive(ref remoteIPE);     // UDP接收数据  
                            if (data_recv.Length > 0 && remoteIPE.Address.ToString() == ServerIPE.Address.ToString()) // 只处理特定的服务端的数据  
                            {
                                string cmdstr = Encoding.UTF8.GetString(data_recv); // 接收命令数据
                                CmdMessage msg = JsonHelper.Deserialize<CmdMessage>(cmdstr);
                                switch (msg.CmdType)
                                {
                                    case "VideoCaptureCommand":
                                        sw.WriteLine("VideoCaptureCommand:  " + msg.JsonText);
                                        VideoCaptureCommand cmd = JsonHelper.Deserialize<VideoCaptureCommand>(msg.JsonText);
                                        if (cmd.CmdType == 1)   // BeginVideo
                                        {
                                            Thread tBegin = new Thread(() =>
                                            {
                                                string dir2 = Properties.Settings.Default.HttpLocalPath;
                                                string file2 = cmd.VideoFile;
                                                string path2 = System.IO.Path.Combine(dir2, file2);

                                                if (!Directory.Exists(dir2)) { Directory.CreateDirectory(dir2); }
                                                string ifmt = (cmd.CmdID == "Screen") ? "-f gdigrab" : "";
                                                string iduration = (cmd.VideoDuration > 0) ? string.Format("-t {0}", cmd.VideoDuration) : "";
                                                string input = (cmd.CmdID == "Screen") ? "desktop" : string.Format("\"rtsp://{0}:{1}@{2}:{3}/cam/realmonitor?channel=1&subtype=0\"", cmd.TargetUser, cmd.TargetPass, cmd.TargetIP, cmd.TargetPort);
                                                string outputArgs = "-vcodec libx264 -pix_fmt yuv420p -acodec ac3";
                                                string output = string.Format("\"{0}\"", path2);
                                                // ffmpeg -t 120 -i "rtsp://admin:xywadmin@192.168.0.151:554/cam/realmonitor?channel=1&subtype=0" -vcodec libx264 -pix_fmt yuv420p -acodec aac "D:\Workspace\web\nginx_v1.7.11.3_Gryphon\nginx-rtmp-module\tmp\rec\y.mp4"
                                                string args = string.Format(@"{0} {1} -i {2} {3} {4}", ifmt, iduration, input, outputArgs, output);
                                                this.cmd.Execute("ffmpeg.exe ", args);

                                            });

                                            tBegin.Start();
                                        }
                                        else if (cmd.CmdType == 2)  // StopVideo
                                        {
                                            Thread tStop = new Thread(new ThreadStart(delegate()
                                            {
                                                if (this.cmd.Running()) { this.cmd.Stop(); }
                                            }));

                                            tStop.Start();
                                        }
                                        else // PauseVideo
                                        {

                                        }
                                        break;
                                    case "CommandSendlog":
                                        sw.WriteLine("CommandSendlog:  " + msg.JsonText);
                                        CommandSendlog cmdlog = JsonHelper.Deserialize<CommandSendlog>(msg.JsonText);
                                        break;
                                    default:
                                        sw.WriteLine(cmdstr);
                                        break;
                                }
                            }
                        }
                        catch (Exception ex)
                        {
                            sw.WriteLine(ex.StackTrace);
                        }
                    }

                    Thread.Sleep(50);
                }
            }
        }

        public void Dispose()
        {
            this.cmd.Stop();
        }
    }

    

    



}