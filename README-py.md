# [StaffSearch](https://staffsearch.fernuni-hagen.de) consultation –  staffs.py

## Prerequisites

You will need a decently recent Python&nbsp;3 installation, probably Python&nbsp;3.8+. You should have *venv* Python module installed and ready to use.

Apart from a standard python installation and its standard library this script makes use of the HTML parser *BeautifulSoup*. Thats all.

## Installation

1. Create a folder, `cd` into it and copy `staff.py` there.
   ```bash
   $ mkdir staff && cd staff
   $ cp source_path/staff.py ./
   ```
2. Create a Python *venv* (you may have to use `python3`/`pip3` instead of `python`/`pip`depending on your OS distribution!):
   ```bash
   $ python -m venv ./venv
   ```
4. Activate the *virtual environment* (you will have to have this *venv* active each time you use **staff.py**; you can choose to install the dependency below *globally*, if you want to avoid that):
   ```bash
   $ source venv/bin/activate
   ```
5. Install HTML parsing dependency:
   ```bash
   $ pip install beautifulsoup4
   ```
6. Optionally make `staff.py` executable:
   ```bash
   $ chmod 0750 staff.py
   ```
7. To leave the *virtual environment* execute the shell function
   ```bash
   $ deactivate
   ```


## Running the command

1. Change into respective directory.
   Activate the *virtual environment* (you will have to have this *venv* active each time you use **staff.py**; you can choose to install the dependency below *globally*, if you want to avoid that):
   ```bash
   $ source venv/bin/activate
   ```
2. Get help on the available command line options:
   ```bash
   $ python staff.py --help
   ```
   Or search for a staff member:
   ```bash
   $ python staff.py ‹search_terms›
   ```
4. To leave the *virtual environment* execute the shell function
   ```bash
   $ deactivate
   ```
