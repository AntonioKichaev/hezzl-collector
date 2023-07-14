CREATE TABLE IF NOT EXISTS event.items(
    Id int,
    CampaignId int,
    Name VARCHAR(255),
    Description VARCHAR(255),
    Priority int,
    Removed bool,
    EventTime datetime
) ENGINE=MergeTree()
order by (Id,CampaignId,Name);
insert into event.items(Id, CampaignId, Name, Description, Priority, Removed, EventTime)
values (55,55,'item 55','',1,false,now());
