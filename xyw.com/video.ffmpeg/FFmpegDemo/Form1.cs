using System;
using System.Diagnostics;
using System.Threading;
using System.Windows.Forms;
using FFmpegService;
using System.IO;


namespace FFmpegDemo
{
    public partial class Form1 : Form
    {
        FFmpegCmd cmd = new FFmpegCmd();
        Process proc = null;
        Thread th = null;

        public Form1()
        {
            InitializeComponent();
        }

        private void button_Start_Click(object sender, EventArgs e)
        {
            // ffmpeg -t 120 -i "rtsp://admin:xywadmin@192.168.0.151:554/cam/realmonitor?channel=1&subtype=0" -vcodec libx264 -pix_fmt yuv420p -acodec aac "D:\Workspace\web\nginx_v1.7.11.3_Gryphon\nginx-rtmp-module\tmp\rec\y.mp4"
            // ffmpeg -f gdigrab -t 120 -i desktop -vcodec libx264 -pix_fmt yuv420p -acodec aac xxxyyy.mp4

            string ffmpeg = @"ffmpeg.exe";
            string file = DateTime.Now.ToString("yyyyMMdd_HHmmss");
            if (File.Exists(file)) { file = DateTime.Now.ToString("yyyyMMdd_HHmmss.fff"); }
            string args = string.Format(@" -f gdigrab -t 120 -i desktop -vcodec libx264 -pix_fmt yuv420p -acodec aac {0}.mp4", file);
            System.Console.WriteLine(ffmpeg + args);
            // this.cmd.Execute(exe, args);
            // this.proc = this.cmd.RunProgram2(ffmpeg, "", args);  

            this.proc = new Process();
            this.proc.StartInfo.UseShellExecute = false;          // 不使用系统外壳程序启动  
            this.proc.StartInfo.RedirectStandardInput = true;    // 不重定向输入  
            this.proc.StartInfo.RedirectStandardOutput = true;    // 重定向输出  
            this.proc.StartInfo.RedirectStandardError = true;
            this.proc.StartInfo.CreateNoWindow = true;            // 不创建窗口  
            this.proc.StartInfo.WindowStyle = ProcessWindowStyle.Hidden;
            this.proc.StartInfo.FileName = ffmpeg;
            this.proc.StartInfo.Arguments = args;
            this.proc.Start();

            new Thread(() =>
            {
                using (System.IO.StreamWriter tssw = new System.IO.StreamWriter(file + ".log", true))
                {
                    tssw.WriteLine(this.proc.StandardError.ReadToEnd());
                }
            }).Start();


            this.button_Start.Enabled = !this.button_Start.Enabled;
            this.button_Stop.Enabled = !this.button_Start.Enabled;
        }

        private void button_Stop_Click(object sender, EventArgs e)
        {
            if (this.proc != null && (!this.proc.HasExited))
            {
                this.proc.StandardInput.WriteLine("{0}", "q");
                this.proc.Dispose();
                this.proc = null;
            }

            this.button_Stop.Enabled = !this.button_Stop.Enabled;
            this.button_Start.Enabled = !this.button_Stop.Enabled;
        }

        string GetFilename()
        {
            string file = DateTime.Now.ToString("yyyyMMdd_HHmmss");
            if (File.Exists(file))
            {
                file = DateTime.Now.ToString("yyyyMMdd_HHmmss.fff");
            }

            return file;
        }
    }
}
