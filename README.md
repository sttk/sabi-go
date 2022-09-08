# Sabi

A small framework for Golang applications.

## What is this?

Sabi is a small framework to separate logics and data accesses for Golang applications.

A program consists of procedures and data.
And to operate data, procedures includes data accesses, then the rest of procedures except data accesses are logics.
Therefore, a program consists of logics, data accesses and data.

This package is an application framework which explicitly separates procedures into logics and data accesses as layers.
By using this framework, we can remove codes for data accesses from logic parts, and write only specific codes for each data source (e.g. database, messaging services files, and so on)  in data access  parts. 
Moreover, by separating as layers, applications using this framework can change data sources easily by switching data access parts.

## License

Copyright (C) 2022 Takayuki Sato

This program is free software under MIT License.<br>
See the file LICENSE in this distribution for more details.
