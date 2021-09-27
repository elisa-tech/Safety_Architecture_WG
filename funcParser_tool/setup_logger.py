"""
SPDX-License-Identifier: GPL-2.0-only
SPDX-FileCopyrightText: Copyright (C) 2021 Intel Corporation

This program is free software; you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free Software
Foundation; version 2.
This program is distributed in the hope that it will be useful, but WITHOUT ANY
WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
PARTICULAR PURPOSE. See the GNU General Public License for more details.
You should have received a copy of the GNU General Public License along with
this program; if not, write to the Free Software Foundation, Inc., 51 Franklin
Street, Fifth Floor, Boston, MA 02110-1301, USA.

Author: Stefano Dell'Osa <stefano.dellosa@intel.com>
"""

import logging

FORMAT = '%(levelname)s:%(message)s'
# Create logger
logging.basicConfig(format=FORMAT)
logger = logging.getLogger('funcParser')
logger.setLevel(logging.INFO)
# Db logger
db_logger = logging.getLogger('sqlalchemy.engine')
db_logger.setLevel(logging.WARNING)
logger.addHandler(db_logger)
