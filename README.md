<h1 align="center">Red-Team BUTLER</h1>

> [!Note]
> Under active development, welcome to submit issues and pull-requests

**"Red-Team BUTLER" (tentative name) is an automated tool that assists Red-Team personnel in the process of performing targets (WEB) penetration testing. Execute automated workflows by customizing workflows and editing custom work methods (plug-ins).**

In future plans, more WMPs (plug-ins) will be developed to cover the regular "reconnaissance" phase of the Red-Team's work. 

Everyone is welcome to provide suggestions, contribute code, and participate in testing.

> [!Warning]
> This tool is intended for use in legitimate penetration testing and security assessments. Before using this tool, make sure you have explicit authorization from the owner of the target system.
> By using this tool, you agree to abide by local laws and bear sole responsibility for any consequences and impacts caused by the use of this tool.

## Summary

The core of BUTLER includes three parts: workflow, wmpci, and grid, which are respectively responsible for workflow execution, work method plugin manage and connection, and planning grid arrangement. 

Create WMP (Work method plugin) according to the requirements of WMPCI (work method plugin common interface) and connect it through the connector. After creating the workflow task, WMP can be scheduled to complete the task in an orderly manner according to the arrangement of Grid.

WMP can be customized to suit individual work methods, and the Grid can be customized to suit individual workflows. 

**Red-Team members can automate their workflows by encapsulating all the tools used in the "reconnaissance" workflow into "WMP" and the workflow itself into "Grid". The combination of "WMP" and "Grid" covers most of their "reconnaissance" workflow needs, and "Workflow" encapsulates the basic attributes of the reconnaissance target into "Asset".**


<hr/>
<p align="center">Thank you for supporting the "Red-Team BUTLER" project</p>
