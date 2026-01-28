package database

import (
	"log"

	"github.com/wayne-stays/backend/internal/models"
	"github.com/wayne-stays/backend/internal/storage"
	"gorm.io/gorm"
)

// Seed populates the database with initial data
func Seed(db *gorm.DB, store storage.Storage) error {
	log.Println("Seeding database...")

	// Check if accommodations already exist
	var count int64
	db.Model(&models.Accommodation{}).Count(&count)
	if count > 0 {
		log.Println("ℹ️  Database already seeded, skipping...")
		return nil
	}

	accommodations := []models.Accommodation{
		{
			Slug:        "wayne-tower-sky-suites",
			Title:       "Wayne Tower Sky Suites",
			Description: "Executive suites in the iconic Wayne Tower. Direct access to Wayne Enterprises board rooms, private viewing galleries, and the legendary penthouse lounge with 360° views from Gotham's tallest building.",
			Location:    "Financial District, Gotham City",
			PriceStart:  185000, // $1,850.00 in cents
			Capacity:    4,
			Category:    "urban-luxury",
			Status:      "active",
			Amenities:   []string{"helipad", "concierge", "skylounge", "board-room-access", "panoramic-views"},
			Images: models.ImageSet{
				Hero: "accommodations/wayne_tower_sky_suites/1.jpg",
				Gallery: []string{
					"accommodations/wayne_tower_sky_suites/2.jpg",
					"accommodations/wayne_tower_sky_suites/3.jpg",
					"accommodations/wayne_tower_sky_suites/4.jpg",
					"accommodations/wayne_tower_sky_suites/5.jpg",
				},
			},
		},
		{
			Slug:        "miagani-island-sanctuary",
			Title:       "Miagani Island Sanctuary",
			Description: "Secluded island retreat accessible only by private ferry. Complete privacy, lighthouse views, and untouched natural beauty. An entire island experience with dedicated staff.",
			Location:    "Miagani Island, Gotham Harbor",
			PriceStart:  320000, // $3,200.00
			Capacity:    8,
			Category:    "island",
			Status:      "active",
			Amenities:   []string{"private-island", "ferry-service", "lighthouse", "dedicated-staff", "beach-access"},
			Images: models.ImageSet{
				Hero: "accommodations/miagani_island_sanctuary/496521281_18341349961091194_7188144892539395392_n.jpg",
				Gallery: []string{
					"accommodations/miagani_island_sanctuary/496518690_18341349970091194_4571822253732037825_n.jpg",
					"accommodations/miagani_island_sanctuary/496776631_18341350006091194_555185071353290621_n.jpg",
					"accommodations/miagani_island_sanctuary/496784740_18341349988091194_3745965920143058440_n.jpg",
					"accommodations/miagani_island_sanctuary/496935162_18341349997091194_7526324161251226011_n.jpg",
					"accommodations/miagani_island_sanctuary/496996305_18341349979091194_2081130433048635428_n.jpg",
					"accommodations/miagani_island_sanctuary/497175762_18341350015091194_3481747048423755075_n.jpg",
					"accommodations/miagani_island_sanctuary/497746403_18341350024091194_471665051166189818_n.jpg",
				},
			},
		},
		{
			Slug:        "clocktower-penthouse",
			Title:       "The Clocktower Penthouse",
			Description: "Luxury apartment within Gotham's famous clocktower. Vaulted ceilings, original clockwork machinery, and stunning Gothic architecture. Live inside a functioning historic landmark.",
			Location:    "Burnley, Old Gotham",
			PriceStart:  140000, // $1,400.00
			Capacity:    2,
			Category:    "historic",
			Status:      "active",
			Amenities:   []string{"historic-building", "vaulted-ceilings", "clockwork-machinery", "gothic-architecture", "city-views"},
			Images: models.ImageSet{
				Hero: "accommodations/clocktower_penthouse/574154632_18533050843043506_4346848479874717654_n.jpg",
				Gallery: []string{
					"accommodations/clocktower_penthouse/571363933_18533050795043506_5440752749890542567_n.jpg",
					"accommodations/clocktower_penthouse/571737091_18533050759043506_2738501978407641211_n.jpg",
					"accommodations/clocktower_penthouse/572111397_18533050831043506_1305076926679561863_n.jpg",
					"accommodations/clocktower_penthouse/572178603_18533050792043506_234395849699942906_n.jpg",
					"accommodations/clocktower_penthouse/572561856_18533050822043506_1111370863296291170_n.jpg",
					"accommodations/clocktower_penthouse/574041038_18533050777043506_1823204043955656439_n.jpg",
					"accommodations/clocktower_penthouse/574238279_18533050810043506_2689492429418356971_n.jpg",
					"accommodations/clocktower_penthouse/574344230_18533050813043506_8599427814650356724_n.jpg",
					"accommodations/clocktower_penthouse/574391371_18533050768043506_8876097839108080990_n.jpg",
				},
			},
		},
		{
			Slug:        "arkham-district-gentlemens-club",
			Title:       "Arkham District Gentleman's Club",
			Description: "Restored Victorian gentleman's club with leather libraries, billiard rooms, and private dining. A step back to Gotham's gilded age with members-only access to exclusive lounges.",
			Location:    "Arkham, North Gotham",
			PriceStart:  95000, // $950.00
			Capacity:    2,
			Category:    "historic",
			Status:      "active",
			Amenities:   []string{"library", "billiards", "private-dining", "victorian-architecture", "exclusive-lounge"},
			Images: models.ImageSet{
				Hero: "accommodations/arkham_district_gentlemens_club/587092492_18332589751210508_2764734941953359835_n.heic.jpg",
				Gallery: []string{
					"accommodations/arkham_district_gentlemens_club/588043143_18332589760210508_572137201056803845_n.heic.jpg",
					"accommodations/arkham_district_gentlemens_club/589122197_18332589742210508_2115549378885077875_n.heic.jpg",
					"accommodations/arkham_district_gentlemens_club/589815981_18332589733210508_4757641263690381477_n.heic.jpg",
				},
			},
		},
		{
			Slug:        "kane-county-vineyard-estate",
			Title:       "Kane County Vineyard Estate",
			Description: "Rolling vineyard estate producing Gotham's finest wines. Private wine cellar, tasting room, and farm-to-table cuisine by Michelin-starred chefs with exclusive vintage tastings.",
			Location:    "Kane County, Gotham Outskirts",
			PriceStart:  210000, // $2,100.00
			Capacity:    6,
			Category:    "countryside",
			Status:      "active",
			Amenities:   []string{"vineyard", "wine-tasting", "fine-dining", "wine-cellar", "countryside"},
			Images: models.ImageSet{
				Hero: "accommodations/kane_county_vineyard_estate/588486778_18535526878045145_7335218717638223358_n.jpg",
				Gallery: []string{
					"accommodations/kane_county_vineyard_estate/587996475_18535526887045145_7036026020694177932_n.jpg",
					"accommodations/kane_county_vineyard_estate/588396384_18535526896045145_2912282764379893080_n.jpg",
					"accommodations/kane_county_vineyard_estate/588715609_18535526869045145_8002033117224592767_n.jpg",
					"accommodations/kane_county_vineyard_estate/588806137_18535526905045145_1387029417042462886_n.jpg",
				},
			},
		},
	}

	// Create all accommodations
	for _, acc := range accommodations {
		if err := db.Create(&acc).Error; err != nil {
			log.Printf("❌ Failed to seed accommodation: %s - %v", acc.Title, err)
			return err
		}
		log.Printf("✅ Seeded: %s", acc.Title)
	}

	// Seed Merchandise
	merchandise := []models.Merchandise{
		{Name: "Wayne Enterprises Leather Portfolio", Price: 32000, Category: "Accessories", IsNew: true, ImageURL: "merch/accessories/portfolio.png"},
		{Name: "Wayne Manor Heritage Cashmere Scarf", Price: 45000, Category: "Apparel", IsNew: true, ImageURL: "merch/apparel/scarf.png"},
		{Name: "Martha's Pearls (Replica)", Price: 45000, Category: "Accessories", ImageURL: "merch/accessories/pearls.png"},
		{Name: "Wayne Manor: Architectural Legacy", Price: 12500, Category: "Home", ImageURL: "merch/home/book.png"},
		{Name: "Martha Wayne Memorial Rose Set", Price: 9500, Category: "Home", ImageURL: "merch/home/roses.png"},
		{Name: "Wayne Foundation Gala Program", Price: 35000, Category: "Collectibles", ImageURL: "merch/collectibles/program.png"},
		{Name: "Kane County Reserve Wine", Price: 42000, Category: "Lifestyle", ImageURL: "merch/lifestyle/wine.png"},
		{Name: "Gotham Skyline Silk Tie", Price: 18500, Category: "Apparel", ImageURL: "merch/apparel/tie.png"},
		{Name: "Wayne Stays Luxury Linen Set", Price: 38000, Category: "Home", ImageURL: "merch/home/linens.png"},
		{Name: "Gotham Opera House Poster", Price: 16500, Category: "Collectibles", ImageURL: "merch/collectibles/poster.png"},
	}

	for _, item := range merchandise {
		if err := db.Create(&item).Error; err != nil {
			log.Printf("❌ Failed to seed merchandise: %s - %v", item.Name, err)
			// Don't fail the whole seed if dupes (though check above is for accommodations only, so this might be needed)
			continue
		}
		log.Printf("✅ Seeded Item: %s", item.Name)
	}

	log.Println("✅ Database seeded successfully")
	return nil
}
