package hw10programoptimization

import (
	"bytes"
	"testing"
)

var data = `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov", "Phone":"6-866-899-36-79", "Password":"InAQJvsq","Address":"Blackbird Place 25"},
  {"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00", "Password":"SiZLeNSGn","Address":"Fulton Hill 80"},
  {"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97", "Password":"71kuz3gA5w","Address":"Monterey Park 39"},
  {"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16", "Password":"r639qLNu","Address":"Sunfield Park 20"},
  {"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01", "Password":"acSBF5","Address":"Russell Trail 61"},
  {"Id":6,"Name":"Test Horse","Username":"Horsie123","Email":"horsie@top.top.top.com","Phone":"146-91-01", "Password":"acSBF5","Address":"Russell Trail 61"},
  {"Id":7,"Name":"Alice Cooper","Username":"AliceC99","Email":"alice@rockstar.net","Phone":"555-867-5309", "Password":"RockNroll","Address":"Sunset Blvd 100"},
  {"Id":8,"Name":"Bob Dylan","Username":"FolkLegend","Email":"bob@music.com","Phone":"555-432-1234", "Password":"TambourineMan","Address":"Highway 61"},
  {"Id":9,"Name":"Charlie Brown","Username":"GoodGrief","Email":"charlie@peanuts.com","Phone":"555-222-3333", "Password":"Snoopy123","Address":"Peanuts Lane 5"},
  {"Id":10,"Name":"Dana White","Username":"UFCBoss","Email":"dana@ufc.com","Phone":"555-999-0000", "Password":"FightNight","Address":"Octagon Drive 88"},
  {"Id":11,"Name":"Eve Adams","Username":"EveA","Email":"eve@genesis.com","Phone":"555-111-2233", "Password":"ForbiddenFruit","Address":"Eden Garden 7"},
  {"Id":12,"Name":"Frank Castle","Username":"Punisher","Email":"frank@justice.com","Phone":"555-666-7777", "Password":"Vigilante","Address":"Hell's Kitchen 13"},
  {"Id":13,"Name":"Grace Kelly","Username":"PrincessGrace","Email":"grace@monaco.mc","Phone":"555-888-1234", "Password":"RoyalLife","Address":"Palace Street 1"},
  {"Id":14,"Name":"Hank Moody","Username":"WriterHank","Email":"hank@californication.tv","Phone":"555-432-5678", "Password":"WhiskeyAndWords","Address":"Venice Beach 33"},
  {"Id":15,"Name":"Isaac Newton","Username":"GravityGuy","Email":"newton@physics.org","Phone":"555-000-9999", "Password":"AppleFall","Address":"Cambridge Univ 10"},
  {"Id":16,"Name":"John Lennon","Username":"Imagine","Email":"john@beatles.com","Phone":"555-654-3210", "Password":"GivePeace","Address":"Abbey Road 42"},
  {"Id":17,"Name":"Kurt Cobain","Username":"TeenSpirit","Email":"kurt@nirvana.com","Phone":"555-777-1111", "Password":"SmellsLike","Address":"Seattle Grunge 9"},
  {"Id":18,"Name":"Leonardo Da Vinci","Username":"RenaissanceMan","Email":"leo@art.com","Phone":"555-555-5555", "Password":"MonaLisa","Address":"Florence Italy 14"},
  {"Id":19,"Name":"Marie Curie","Username":"RadiumQueen","Email":"marie@science.org","Phone":"555-333-2222", "Password":"Radioactive","Address":"Paris Science Lab 50"},
  {"Id":20,"Name":"Nikola Tesla","Username":"ElectricMaster","Email":"nikola@tesla.com","Phone":"555-999-8888", "Password":"ACDC","Address":"Wardenclyffe Tower 1"},
  {"Id":21,"Name":"Oscar Wilde","Username":"WittyWriter","Email":"oscar@literature.com","Phone":"555-444-7777", "Password":"TheImportance","Address":"Dublin Writers Row 3"},
  {"Id":22,"Name":"Pablo Picasso","Username":"Cubist","Email":"pablo@art.com","Phone":"555-222-6666", "Password":"BluePeriod","Address":"Barcelona Art St 8"},
  {"Id":23,"Name":"Quentin Tarantino","Username":"QTDirector","Email":"quentin@movies.com","Phone":"555-121-2121", "Password":"PulpFiction","Address":"Hollywood Blvd 90"},
  {"Id":24,"Name":"Robin Hood","Username":"ArcherKing","Email":"robin@nottingham.uk","Phone":"555-343-2323", "Password":"Sherwood","Address":"Sherwood Forest 77"},
  {"Id":25,"Name":"Stephen Hawking","Username":"CosmosThinker","Email":"stephen@theory.com","Phone":"555-999-2222", "Password":"BlackHoles","Address":"Cambridge Univ 20"},
  {"Id":26,"Name":"Tom Sawyer","Username":"FencePainter","Email":"tom@marktwain.com","Phone":"555-565-6767", "Password":"RiverAdventures","Address":"Mississippi River 4"},
  {"Id":27,"Name":"Usain Bolt","Username":"FastestMan","Email":"usain@speed.com","Phone":"555-909-1010", "Password":"Lightning","Address":"Kingston Jamaica 6"},
  {"Id":28,"Name":"Vincent Van Gogh","Username":"StarryPainter","Email":"vincent@art.com","Phone":"555-777-9876", "Password":"Sunflowers","Address":"Arles France 22"}`

func BenchmarkDomainStat(b *testing.B) {
	bufferedData := bytes.NewBufferString(data)
	for i := 0; i < b.N; i++ {
		stat, err := GetDomainStat(bufferedData, "com")
		_, _ = stat, err
	}
}
