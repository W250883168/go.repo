using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FFmpegService
{
    public class CmdMessage
    {
        public string CmdID { get; set; }
        public string CmdType { get; set; }
        public string JsonText { get; set; }

        public CmdMessage()
        {
            this.CmdID = "";
            this.CmdType = "";
            this.JsonText = "";
        }
    }
}
