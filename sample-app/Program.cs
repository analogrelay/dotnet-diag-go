using System;
using System.Threading.Tasks;

namespace SampleApp
{
    class Program
    {
        static async Task Main(string[] args)
        {
            var tcs = new TaskCompletionSource<object>();
            Console.CancelKeyPress += (sender, a) => {
                a.Cancel = true;
                Console.WriteLine("Stopping...");
                tcs.TrySetResult(null);
            };
            Console.WriteLine("Running. Press Ctrl-C to stop.");
            await tcs.Task;
        }
    }
}
