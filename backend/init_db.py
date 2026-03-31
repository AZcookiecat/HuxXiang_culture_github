from app import create_app, db
from models.cultural_resource import CulturalResource
from models.user import User


def init_database():
    app = create_app()

    with app.app_context():
        print("Creating Flask backend tables...")
        db.create_all()
        print("Flask backend tables created.")

        admin = User.query.filter_by(username="admin").first()
        if not admin:
            print("Creating default admin user...")
            admin = User(
                username="admin",
                email="admin@example.com",
                role="admin",
            )
            admin.set_password("admin123")
            db.session.add(admin)
            db.session.commit()
            print("Default admin user created.")
        else:
            print("Default admin user already exists.")

        sample_resource = CulturalResource.query.filter_by(title="Huxiang culture overview").first()
        if not sample_resource:
            print("Creating sample cultural resource...")
            sample_resource = CulturalResource(
                title="Huxiang culture overview",
                description="Sample cultural resource used to verify the Flask resource module.",
                content="Huxiang culture is the regional culture of Hunan and remains available in the Flask backend after the post system was moved to Gin.",
                type="history",
                category="introduction",
                tags="huxiang,hunan,culture,history",
                author="system",
                status="published",
                priority=1,
            )
            db.session.add(sample_resource)
            db.session.commit()
            print("Sample cultural resource created.")
        else:
            print("Sample cultural resource already exists.")

        print("Flask backend initialization finished.")


if __name__ == "__main__":
    init_database()
