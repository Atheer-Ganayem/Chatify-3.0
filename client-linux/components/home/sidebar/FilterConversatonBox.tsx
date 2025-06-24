import { Search } from "lucide-react";

import { Label } from "@/components/ui/label";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarInput,
} from "@/components/ui/sidebar";
import AddConversation from "./AddConversationDialog";

interface Props {
  currentFilter: string;
  onFilter: (str: string) => void;
}

const FilterConversatonBox: React.FC<Props> = ({ onFilter, currentFilter }) => {
  return (
    <div className="flex">
      <div className="w-full">
        <SidebarGroup className="py-0">
          <SidebarGroupContent className="relative">
            <Label htmlFor="search" className="sr-only">
              Search
            </Label>
            <SidebarInput
              placeholder="Search conversation..."
              className="pl-8"
              value={currentFilter}
              onChange={(e) => onFilter(e.target.value)}
            />
            <Search className="pointer-events-none absolute top-1/2 left-2 size-4 -translate-y-1/2 opacity-50 select-none" />
          </SidebarGroupContent>
        </SidebarGroup>
      </div>
      <AddConversation />
    </div>
  );
};

export default FilterConversatonBox;
