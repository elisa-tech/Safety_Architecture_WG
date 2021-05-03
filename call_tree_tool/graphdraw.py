# Copyright (C) 2021 Intel Corporation
# SPDX-License-Identifier: GPL-2.0-only

'''
Graph Draw
Draws the call graph for a given function
'''

class GraphDrawFile:
    '''Draw a call tree.'''
    # origin is (shortpath, name, coverage)
    def __init__(self, database, origin, outfilename, depth):
        self.database = database
        self.origin = origin
        self.outfilename = outfilename
        self.indentation = 0
        # stack is (item, depth)
        self.stack = [origin, 0]
        self.depth = depth

    def dfs(self, outfile):
        '''Draw a call tree by performing a depth first search.'''
        # stack is (origin, depth)
        stack = [(self.origin, 0)]
        while len(stack) > 0:
            node = stack.pop(0)
            element = node[0]
            depth = node[1]
            funcname = element[1]
            shortpath = element[0]
            entry = '\t' * depth + funcname + " (" + shortpath + ")"
            print(entry)
            outfile.write(entry + "\n")
            if depth <= self.depth:
                children = self.database.get_edges(element)
                for child in children:
                    stack.insert(0, (child, depth + 1))

    def draw(self):
        '''Draw a call tree and store it in a file.'''
        with open(self.outfilename, "w") as outfile:
            self.dfs(outfile)
