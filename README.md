# Million Video Builder

This program is used to generate a video sequence of the digits 1-1,000,000 based on a series of input video clips

## Requirements

1. Windows 10 OS
1. 32GB Ram (will work with less RAM, but will be much slower)
1. Make sure you have FFmpeg installed and accessible from command line
1. Folder containing `.MP4` video clips for digits 0-9 with the naming convention `#.mp4`, as well as a single clip titled `gap.mp4`

## How to run

1. Place `million-video-builder.exe` into the root project directory
1. With the command line pointed to the project directory, run `.\million-video-builder.exe [src_folder\path] [path\to\output_folder]` (without [ ] braces)

## Pausing, Resuming and Restarting

If you want to pause the render, you can do so easily by exiting the command prompt window. The render will resume where it left off next time the program is run (as long as the output directory is the same)

If you would like to restart the process from scratch, simply delete the existing output directory or provide a new output directory to render to.