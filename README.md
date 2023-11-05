# Sword of Inevitable Doom

It will be a generic roguelike, with an ECS of some description. Trying to hand
roll as much as possible except for all the graphics stuff, using ebitengine
under the hood.

I really don't know how any of this will pan out but its been fun designing the
ECS - though I know for sure it won't perform well for games that have a lot of
objects, I'm taking the approach that I will optimize later, just getting it to
work is the first step.

## TODO

-   Some kind of event system, so when something happens, other systems can
    react to things. Currently I'm just putting arrays on components to record
    changes being made from systems, maybe that will work. I dunno.

# License

MIT License

Copyright (c) 2023 Nathan Ollerenshaw

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
