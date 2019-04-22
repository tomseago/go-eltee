Glossary
========

**Profile** - A class of fixture. Think an individual model from a particular manufacturer.

**Profile Control** - A capability of the lighting fixture. A way in which it can change. These can be grouped together. The color of the light is a good example of a control. The direction in which the light is pointing is another control. 

**Fixture** - A real actual light. As an object in ElTee this is an instantiation of a specific profile assigned to a specific DMX base address in a universe

**Fixture Control** - An instance of a Profile Control. Fixture Controls observe Control Points in the World State through their Lens Stack.

**Lens Stack** - An ordered collection of lenses

**Lens** - An object which modifies the value of an observed Control Point. Lenses are used to provide corrections such as color correction or spatial positioning differences such as parallax. Lenses are the last operators in the control point domain before the fixture control adapts the control point into DMX output.

**Control Point** - A typed value, usually multidimensional, which exists within the control domain. These represent generic concepts such as "the primary color" or "the primary aim point" which zero or more fixtures controls may be observing. Control Points are observed by Fixture Controls and they may be modified by Animators or input adapters.

**World State** - A set of named and typed control points, like colors and positions, which change over time. Animations manipulate the world state, but so do ad hoc human inputs. At one moment in time there is only one world state, but future world states may be known. They just aren't THE state.

**Key State** - A partial world state which can be transitioned (tween'ed) from or to. Control Points are essentially always either sitting at a key state or are tweening between two key states using a tweening algorithm.

**Tweening algorithm** - A way to move between two values over time. Since the values contained in the world state are typed (color, position, etc) we can speak of tweening algorithms generally but the details of how they work are type specific.

**Fixture Patch** - A particular set of connections between fixture controls and control points. Each fixture control can observe only a single control point, but multiple fixture controls can observe the same control point. Fixture Patches can be transitioned between each other either as a hard cut or as a fade possessing some duration. If faded, internally the system will create a temporary version of the arrival patch which matches the values of the departure patch, but has the fixtures attached to the new control points. It will then animate the control points from the departure to the arrival values following whatever easing curves may be specified.

**Controller Patch** - A particular set of connections between controller instances and the control points. Unlike fixture patches, controller patches are two way in that controller instances will be notified as control point values they are connected to change. 

Object Relationships
====================

There is a library of **Profiles** which define types of fixtures. So when a new type of light is introduced into the rig, a new **Profile** needs to be created for it before ElTee can do anything with it. This means writing a new configuration file which will define **Profile Controls** of different types and how those map to DMX address offsets or whatever else will be needed to control that type of light. Each **Profile Control** has both an Id and a Name. The Ids must be unique.

The next set of definitions is for **Fixtures** which are specific instances of **Profiles**. In the land of DMX this means associating a **Profile** by name with a specific DMX base address. **Fixtures** also include a name like "Flood Par 1" that is used to refer to that specific fixture in the future. When **Fixtures** are instantiated, a set of **Fixture Controls** are also instantiated based on the **Profile Controls** defined for the **Profile**.  These **Fixture Controls** inherit the same id values as the **Profile Controls** they are derived from.

The third set of objects is the set of **World States** that the **State Juggler** uses to know what **Control Points** exist and how to map between them and **Fixtures**. In the configuration directory these are known as deployments. 

Each **World State** consists of a list of named **Control Points** with values and possibly a list of **Fixture Patches** which map between these **Control Points** and specific **Fixtures** by defining a **Lens Stack** that is used by the **Fixture** to observe it's control points. During the generation of each frame of output data, each **Fixture** is told to make this observation, which results in the **Fixture** setting specific DMX values into the DMX output buffer.

The **State Juggler** manages the list of all known **World States** which includes those defined in the configuration files as well as a few states that are created dynamically. The State Juggler maintains an authoritative _current_ state with a list of **Control Points** that the **Fixture Controls** directly hold references to. This current state is built from the base **World State** defined for the deployment.

Initially, the **State Juggler** does not know how to merge one state into another. Instead, it only knows how to "cut" to a new state. It does this by iterating through all **Control Points** defined in the new state and directly assigns their values to the existing **Control Points** in the current state. This means that if a **Control Point** isn't defined in the base state, a future state won't be able to assign that **Control Point** a value.

In addition to applying an entire state, the **Control Points** in the current state can be directly manipulated via the api. This is a low level interface which should be done via more abstract Controller Adapters and Controller Patches perhaps, but we're also not there yet.

In the future, we want the **State Juggler** to gain more interesting audio mixer like abilities. The idea is that you would have multiple state channels, each of which could then have effects or filters applied to it including an intensity modifier which would determine the "volume" of that particular state in terms of how much that state effects the current state. One can additionally imagine animations that apply to states which could be placed on these inputs.

Lots of things will probably happen through the combination of local fixture state and lenses in the lens stack. Instead of layers of unary operations that might then rely on local state we might imagine using multi valued functional operations. The same things can probably be done with both approaches so it will have to be determined in practice what makes more sense. For now, we just do direct assignment from a single control point per control.


Main Looper
===========

The main looper is the central animation and event loop. It goes through the following stages:

1. Polling animators for control point updates
2. Polling input adapters for control point changes
3. Publishing changed control point values to subscribed input adaptors
4. Publishing changed control point values to observing fixture controls
5. Poking all Fixtures to have them generate their DMX output
6. Sending the DMX frame
7. Checking for shutdown/reset conditions???


TODOs
=====
* Complete the update frame function in the server.go file. This is basically the main looper described above
* Implement some input adapters
* Implement animations of control points
* Add a web server for hosting UI
* UI to change control point to fixture control mappings (Fixture Patches)
* UI for Input Patches
* UI For direct DMX control of a fixture
* UI/wizardry for creating a set of control points that directly map one to one to a fixtures controls

Possibly want to add an optimization to FixtureControl updaters where they use the WasDirty() method from their attached control points. However since they might be observing the control point through a lens stack that stack also probably requires a dirty marker so for now we just always recalculate everything from scratch on each frame.

Installation
============

Raspberry PI, starting with Raspian base image.

Get the latest go tarball for ARM from https://golang.org/dl/

Unpack that somewhere and set GOROOT to the location. Add $GOROOT/bin to path.

Create a `go` directory at `~/go` and set GOPATH to this.

Install libftdi1 using
    
    sudo apt-get install libftdi1


Interface
=========

The ElTee server exposes a gRPC interface for all it's fun functionality. This includes management sort of things like learning the configuration or modifying it, as well as modifying control points. For web clients an external gRPC-web proxy (specifically Envoy) is needed in order to convert gRPC-web traffic into gRPC traffic.

The idea is to keep everything gRPC, but if necessary we can always expose something else. For instance there is a websockets interface that was started, but is deprecated in favor of gRPC.



